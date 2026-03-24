/**
 * Lightweight EXIF metadata parser.
 * Reads the EXIF IFD0 tags from a JPEG file to detect whether
 * the image was captured by a real camera/device.
 *
 * Returns { hasExif, tags } where tags includes Make, Model,
 * DateTime, Software, GPS, etc. when available.
 */

const TAG_NAMES = {
  0x010f: "Make",
  0x0110: "Model",
  0x0112: "Orientation",
  0x011a: "XResolution",
  0x011b: "YResolution",
  0x0132: "DateTime",
  0x010e: "ImageDescription",
  0x0131: "Software",
  0x013b: "Artist",
  0x8298: "Copyright",
  0x9003: "DateTimeOriginal",
  0x9004: "DateTimeDigitized",
  0x920a: "FocalLength",
  0xa405: "FocalLengthIn35mmFilm",
  0x829a: "ExposureTime",
  0x829d: "FNumber",
  0x8827: "ISOSpeedRatings",
  0xa001: "ColorSpace",
  0xa002: "PixelXDimension",
  0xa003: "PixelYDimension",
  0x9209: "Flash",
  0xa434: "LensModel",
  0xa433: "LensMake",
  // GPS
  0x0001: "GPSLatitudeRef",
  0x0002: "GPSLatitude",
  0x0003: "GPSLongitudeRef",
  0x0004: "GPSLongitude",
};

function readString(view, offset, length) {
  let str = "";
  for (let i = 0; i < length; i++) {
    const c = view.getUint8(offset + i);
    if (c === 0) break;
    str += String.fromCharCode(c);
  }
  return str.trim();
}

function readTagValue(view, type, count, valueOffset, littleEndian) {
  switch (type) {
    case 2: // ASCII string
      return readString(view, valueOffset, count);
    case 3: // SHORT
      return count === 1 ? view.getUint16(valueOffset, littleEndian) : valueOffset;
    case 4: // LONG
      return count === 1 ? view.getUint32(valueOffset, littleEndian) : valueOffset;
    case 5: // RATIONAL
      return (
        view.getUint32(valueOffset, littleEndian) /
        view.getUint32(valueOffset + 4, littleEndian)
      );
    case 10: // SRATIONAL
      return (
        view.getInt32(valueOffset, littleEndian) /
        view.getInt32(valueOffset + 4, littleEndian)
      );
    default:
      return null;
  }
}

function readIFD(view, tiffStart, ifdOffset, littleEndian) {
  const tags = {};
  const numEntries = view.getUint16(ifdOffset, littleEndian);

  for (let i = 0; i < numEntries; i++) {
    const entryOffset = ifdOffset + 2 + i * 12;
    const tag = view.getUint16(entryOffset, littleEndian);
    const type = view.getUint16(entryOffset + 2, littleEndian);
    const count = view.getUint32(entryOffset + 4, littleEndian);

    // Value size
    const typeSizes = { 1: 1, 2: 1, 3: 2, 4: 4, 5: 8, 7: 1, 9: 4, 10: 8 };
    const totalSize = (typeSizes[type] || 1) * count;

    let valueOffset;
    if (totalSize <= 4) {
      valueOffset = entryOffset + 8;
    } else {
      valueOffset = tiffStart + view.getUint32(entryOffset + 8, littleEndian);
    }

    const tagName = TAG_NAMES[tag];
    if (tagName) {
      try {
        tags[tagName] = readTagValue(view, type, count, valueOffset, littleEndian);
      } catch {
        // skip unreadable tag
      }
    }

    // If this is the EXIF IFD pointer (tag 0x8769), read sub-IFD
    if (tag === 0x8769) {
      const exifOffset = tiffStart + view.getUint32(entryOffset + 8, littleEndian);
      try {
        const subTags = readIFD(view, tiffStart, exifOffset, littleEndian);
        Object.assign(tags, subTags);
      } catch {
        // skip
      }
    }

    // GPS IFD pointer (tag 0x8825)
    if (tag === 0x8825) {
      const gpsOffset = tiffStart + view.getUint32(entryOffset + 8, littleEndian);
      try {
        const gpsTags = readIFD(view, tiffStart, gpsOffset, littleEndian);
        Object.assign(tags, gpsTags);
        tags._hasGPS = true;
      } catch {
        // skip
      }
    }
  }

  return tags;
}

/**
 * Parse EXIF from a File or Blob.
 * @param {File|Blob} file
 * @returns {Promise<{ hasExif: boolean, tags: Record<string, any>, confidence: string }>}
 */
export async function parseExif(file) {
  const noExif = { hasExif: false, tags: {}, confidence: "none" };

  // Only process JPEG/TIFF images
  if (
    file.type &&
    !file.type.includes("jpeg") &&
    !file.type.includes("jpg") &&
    !file.type.includes("tiff")
  ) {
    // PNG and WEBP typically don't carry EXIF — check for PNG specifically
    if (file.type.includes("png") || file.type.includes("webp")) {
      return { ...noExif, confidence: "unknown_format" };
    }
    return noExif;
  }

  try {
    const buffer = await file.arrayBuffer();
    const view = new DataView(buffer);

    // Check JPEG SOI marker
    if (view.getUint16(0) !== 0xffd8) {
      return noExif;
    }

    // Scan for APP1 (EXIF) marker
    let offset = 2;
    while (offset < view.byteLength - 4) {
      const marker = view.getUint16(offset);

      if (marker === 0xffe1) {
        // APP1 marker found
        const length = view.getUint16(offset + 2);

        // Check for "Exif\0\0"
        const exifHeader = readString(view, offset + 4, 4);
        if (exifHeader !== "Exif") {
          offset += 2 + length;
          continue;
        }

        const tiffStart = offset + 10; // After marker (2) + length (2) + "Exif\0\0" (6)
        const byteOrder = view.getUint16(tiffStart);
        const littleEndian = byteOrder === 0x4949; // "II"

        // Verify TIFF magic
        if (view.getUint16(tiffStart + 2, littleEndian) !== 0x002a) {
          return noExif;
        }

        const ifdOffset = tiffStart + view.getUint32(tiffStart + 4, littleEndian);
        const tags = readIFD(view, tiffStart, ifdOffset, littleEndian);

        // Determine confidence level
        let confidence = "low";
        const hasCameraInfo = tags.Make || tags.Model;
        const hasTimestamp = tags.DateTime || tags.DateTimeOriginal;
        const hasLensInfo = tags.LensModel || tags.FocalLength;
        const hasGPS = tags._hasGPS;

        if (hasCameraInfo && hasTimestamp) confidence = "high";
        else if (hasCameraInfo || hasTimestamp) confidence = "medium";
        if (hasGPS) confidence = "high";
        if (hasLensInfo) confidence = "high";

        // Clean up internal flags
        delete tags._hasGPS;

        return {
          hasExif: true,
          tags,
          confidence,
          hasGPS: !!hasGPS,
        };
      }

      // Not APP1, skip to next marker
      if ((marker & 0xff00) !== 0xff00) break;
      const segLength = view.getUint16(offset + 2);
      offset += 2 + segLength;
    }

    return noExif;
  } catch {
    return noExif;
  }
}

/**
 * Quick check: does this image have genuine EXIF camera metadata?
 * @param {File|Blob} file
 * @returns {Promise<boolean>}
 */
export async function hasExifData(file) {
  const result = await parseExif(file);
  return result.hasExif && result.confidence !== "none";
}
