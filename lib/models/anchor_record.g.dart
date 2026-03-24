// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'anchor_record.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class AnchorRecordAdapter extends TypeAdapter<AnchorRecord> {
  @override
  final int typeId = 0;

  @override
  AnchorRecord read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return AnchorRecord(
      id: fields[0] as String,
      sha512Hash: fields[1] as String,
      pHash: fields[2] as String,
      ipfsCID: fields[3] as String,
      txHash: fields[4] as String,
      signerAddress: fields[5] as String,
      merkleRoot: fields[6] as String,
      anchoredAt: fields[7] as DateTime,
      location: fields[8] as String?,
      eventContext: fields[9] as String?,
      author: fields[10] as String?,
      license: fields[11] as String?,
      imagePath: fields[12] as String,
      isOnChain: fields[13] as bool,
      mediaType: fields[14] as String,
    );
  }

  @override
  void write(BinaryWriter writer, AnchorRecord obj) {
    writer
      ..writeByte(15)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.sha512Hash)
      ..writeByte(2)
      ..write(obj.pHash)
      ..writeByte(3)
      ..write(obj.ipfsCID)
      ..writeByte(4)
      ..write(obj.txHash)
      ..writeByte(5)
      ..write(obj.signerAddress)
      ..writeByte(6)
      ..write(obj.merkleRoot)
      ..writeByte(7)
      ..write(obj.anchoredAt)
      ..writeByte(8)
      ..write(obj.location)
      ..writeByte(9)
      ..write(obj.eventContext)
      ..writeByte(10)
      ..write(obj.author)
      ..writeByte(11)
      ..write(obj.license)
      ..writeByte(12)
      ..write(obj.imagePath)
      ..writeByte(13)
      ..write(obj.isOnChain)
      ..writeByte(14)
      ..write(obj.mediaType);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is AnchorRecordAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
