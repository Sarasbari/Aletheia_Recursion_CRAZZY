import 'package:flutter/foundation.dart';

class FraudFlags {
  final bool duplicate;
  final bool spam;
  final bool replay;
  final bool suspicious;

  FraudFlags({
    required this.duplicate,
    required this.spam,
    required this.replay,
    required this.suspicious,
  });

  factory FraudFlags.fromJson(Map<String, dynamic> json) {
    return FraudFlags(
      duplicate: json['duplicate'] ?? false,
      spam: json['spam'] ?? false,
      replay: json['replay'] ?? false,
      suspicious: json['suspicious'] ?? false,
    );
  }
}

class TamperedRegion {
  final int x;
  final int y;

  TamperedRegion({required this.x, required this.y});

  factory TamperedRegion.fromJson(Map<String, dynamic> json) {
    return TamperedRegion(
      x: json['x'] ?? 0,
      y: json['y'] ?? 0,
    );
  }
}

class ProofReportBreakdown {
  final String hash;
  final double similarity;
  final String signature;

  ProofReportBreakdown({
    required this.hash,
    required this.similarity,
    required this.signature,
  });

  factory ProofReportBreakdown.fromJson(Map<String, dynamic> json) {
    return ProofReportBreakdown(
      hash: json['hash'] ?? '',
      similarity: (json['similarity'] ?? 0).toDouble(),
      signature: json['signature'] ?? '',
    );
  }
}

class ProofReport {
  final String proofId;
  final String sha256;
  final String storageRef;
  final String cid;
  final String timestamp;
  final int blockNumber;
  final String verdict;
  final int trustScore;
  final ProofReportBreakdown breakdown;
  final FraudFlags flags;

  ProofReport({
    required this.proofId,
    required this.sha256,
    required this.storageRef,
    required this.cid,
    required this.timestamp,
    required this.blockNumber,
    required this.verdict,
    required this.trustScore,
    required this.breakdown,
    required this.flags,
  });

  factory ProofReport.fromJson(Map<String, dynamic> json) {
    return ProofReport(
      proofId: json['proofId'] ?? '',
      sha256: json['sha256'] ?? '',
      storageRef: json['storageRef'] ?? '',
      cid: json['CID'] ?? '',
      timestamp: json['timestamp'] ?? '',
      blockNumber: json['blockNumber'] ?? 0,
      verdict: json['verdict'] ?? 'Unknown',
      trustScore: json['trustScore'] ?? 0,
      breakdown: ProofReportBreakdown.fromJson(json['breakdown'] ?? {}),
      flags: FraudFlags.fromJson(json['flags'] ?? {}),
    );
  }
}

class AletheiaUploadResponse {
  final String proofId;
  final String sha256;
  final String phash;
  final String merkleRoot;
  final String storageRef;
  final String proofRef;
  final String signature;
  final String publicKey;
  final String captureTimestamp;
  final String ipfsCID;
  final String txHash;
  final int blockNumber;
  final String status;
  final FraudFlags flags;
  final List<String> warnings;

  AletheiaUploadResponse({
    required this.proofId,
    required this.sha256,
    required this.phash,
    required this.merkleRoot,
    required this.storageRef,
    required this.proofRef,
    required this.signature,
    required this.publicKey,
    required this.captureTimestamp,
    required this.ipfsCID,
    required this.txHash,
    required this.blockNumber,
    required this.status,
    required this.flags,
    required this.warnings,
  });

  factory AletheiaUploadResponse.fromJson(Map<String, dynamic> json) {
    return AletheiaUploadResponse(
      proofId: json['proofId'] ?? '',
      sha256: json['sha256'] ?? '',
      phash: json['phash'] ?? '',
      merkleRoot: json['merkleRoot'] ?? '',
      storageRef: json['storageRef'] ?? '',
      proofRef: json['proofRef'] ?? '',
      signature: json['signature'] ?? '',
      publicKey: json['publicKey'] ?? '',
      captureTimestamp: json['captureTimestamp'] ?? '',
      ipfsCID: json['ipfsCID'] ?? '',
      txHash: json['txHash'] ?? '',
      blockNumber: json['blockNumber'] ?? 0,
      status: json['status'] ?? '',
      flags: FraudFlags.fromJson(json['flags'] ?? {}),
      warnings: List<String>.from(json['warnings'] ?? []),
    );
  }
}

class AletheiaVerifyResponse {
  final String verdict;
  final int trustScore;
  final bool timestampValid;
  final String sourceType;
  final ProofReport proofReport;
  final List<String> warnings;
  final bool suspicious;
  final String matchedProofId;
  final int similarityScore;
  final List<TamperedRegion> tamperedRegions;
  final bool replaySuspicious;

  AletheiaVerifyResponse({
    required this.verdict,
    required this.trustScore,
    required this.timestampValid,
    required this.sourceType,
    required this.proofReport,
    required this.warnings,
    required this.suspicious,
    required this.matchedProofId,
    required this.similarityScore,
    required this.tamperedRegions,
    required this.replaySuspicious,
  });

  factory AletheiaVerifyResponse.fromJson(Map<String, dynamic> json) {
    return AletheiaVerifyResponse(
      verdict: json['verdict'] ?? 'Unknown',
      trustScore: json['trustScore'] ?? 0,
      timestampValid: json['timestampValid'] ?? false,
      sourceType: json['sourceType'] ?? '',
      proofReport: ProofReport.fromJson(json['proofReport'] ?? {}),
      warnings: List<String>.from(json['warnings'] ?? []),
      suspicious: json['suspicious'] ?? false,
      matchedProofId: json['matchedProofId'] ?? '',
      similarityScore: json['similarityScore'] ?? 0,
      tamperedRegions: (json['tamperedRegions'] as List?)
              ?.map((i) => TamperedRegion.fromJson(i))
              .toList() ??
          [],
      replaySuspicious: json['replaySuspicious'] ?? false,
    );
  }
}

class VideoUploadResponse {
  final String proofId;
  final String videoHash;
  final String audioHash;
  final String cidVideo;
  final String cidAudio;
  final String txHash;
  final String status;

  VideoUploadResponse({
    required this.proofId,
    required this.videoHash,
    required this.audioHash,
    required this.cidVideo,
    required this.cidAudio,
    required this.txHash,
    required this.status,
  });

  factory VideoUploadResponse.fromJson(Map<String, dynamic> json) {
    return VideoUploadResponse(
      proofId: json['proofId'] ?? '',
      videoHash: json['videoHash'] ?? '',
      audioHash: json['audioHash'] ?? '',
      cidVideo: json['cidVideo'] ?? '',
      cidAudio: json['cidAudio'] ?? '',
      txHash: json['txHash'] ?? '',
      status: json['status'] ?? '',
    );
  }
}

class VideoVerifyDetails {
  final String video;
  final String audio;

  VideoVerifyDetails({required this.video, required this.audio});

  factory VideoVerifyDetails.fromJson(Map<String, dynamic> json) {
    return VideoVerifyDetails(
      video: json['video'] ?? '',
      audio: json['audio'] ?? '',
    );
  }
}

class VideoVerifyResponse {
  final String videoHash;
  final String audioHash;
  final String verdict;
  final VideoVerifyDetails details;

  VideoVerifyResponse({
    required this.videoHash,
    required this.audioHash,
    required this.verdict,
    required this.details,
  });

  factory VideoVerifyResponse.fromJson(Map<String, dynamic> json) {
    return VideoVerifyResponse(
      videoHash: json['videoHash'] ?? '',
      audioHash: json['audioHash'] ?? '',
      verdict: json['verdict'] ?? '',
      details: VideoVerifyDetails.fromJson(json['details'] ?? {}),
    );
  }
}
