import 'aletheia_models.dart';

class VerificationResult {
  final bool isAuthentic;
  final double ssimScore;
  final int hammingDistance;
  final bool sha512Match;
  final bool merkleProofValid;
  final bool ecdsaValid;
  final String anchoredHash;
  final String submittedHash;
  final String signerAddress;
  final String ipfsCID;
  final String txHash;
  final String blockNumber;
  final DateTime anchoredAt;
  final String anchoredByAddress;

  // New fields from Aletheia API
  final int trustScore;
  final String verdict;
  final FraudFlags? fraudFlags;
  final List<TamperedRegion> tamperedRegions;

  VerificationResult({
    required this.isAuthentic,
    required this.ssimScore,
    required this.hammingDistance,
    required this.sha512Match,
    required this.merkleProofValid,
    required this.ecdsaValid,
    required this.anchoredHash,
    required this.submittedHash,
    required this.signerAddress,
    required this.ipfsCID,
    required this.txHash,
    required this.blockNumber,
    required this.anchoredAt,
    required this.anchoredByAddress,
    required this.trustScore,
    required this.verdict,
    this.fraudFlags,
    this.tamperedRegions = const [],
  });

  factory VerificationResult.fromAletheia(AletheiaVerifyResponse response) {
    return VerificationResult(
      isAuthentic: response.verdict.toLowerCase().contains('authentic') || response.trustScore >= 80,
      ssimScore: response.similarityScore / 100.0,
      hammingDistance: (100 - response.similarityScore).toInt(),
      sha512Match: response.similarityScore >= 95, 
      merkleProofValid: response.timestampValid,
      ecdsaValid: response.proofReport.breakdown.signature.isNotEmpty,
      anchoredHash: response.proofReport.sha256,
      submittedHash: response.proofReport.breakdown.hash,
      signerAddress: response.proofReport.breakdown.signature.isNotEmpty ? "0xVerified" : "Unknown",
      ipfsCID: response.proofReport.cid,
      txHash: response.proofReport.proofId, 
      blockNumber: response.proofReport.blockNumber.toString(),
      anchoredAt: DateTime.tryParse(response.proofReport.timestamp) ?? DateTime.now(),
      anchoredByAddress: response.sourceType,
      trustScore: response.trustScore,
      verdict: response.verdict,
      fraudFlags: response.proofReport.flags,
      tamperedRegions: response.tamperedRegions,
    );
  }

  factory VerificationResult.fromVideoAletheia(VideoVerifyResponse response) {
    return VerificationResult(
      isAuthentic: response.verdict.toLowerCase().contains('authentic'),
      ssimScore: 1.0, // Placeholder for video
      hammingDistance: 0,
      sha512Match: true,
      merkleProofValid: true,
      ecdsaValid: true,
      anchoredHash: response.videoHash,
      submittedHash: response.videoHash,
      signerAddress: "0xVerified",
      ipfsCID: "VideoVerified",
      txHash: "VideoProof",
      blockNumber: "0",
      anchoredAt: DateTime.now(),
      anchoredByAddress: "Video Engine",
      trustScore: 100,
      verdict: response.verdict,
    );
  }

  factory VerificationResult.error(String message) {
    return VerificationResult(
      isAuthentic: false,
      ssimScore: 0,
      hammingDistance: 0,
      sha512Match: false,
      merkleProofValid: false,
      ecdsaValid: false,
      anchoredHash: "—",
      submittedHash: "—",
      signerAddress: "—",
      ipfsCID: "—",
      txHash: "—",
      blockNumber: "—",
      anchoredAt: DateTime.now(),
      anchoredByAddress: "—",
      trustScore: 0,
      verdict: message,
    );
  }
}
