import 'package:hive/hive.dart';
import 'package:uuid/uuid.dart';

part 'anchor_record.g.dart';

@HiveType(typeId: 0)
class AnchorRecord extends HiveObject {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String sha512Hash;

  @HiveField(2)
  final String pHash;

  @HiveField(3)
  final String ipfsCID;

  @HiveField(4)
  final String txHash;

  @HiveField(5)
  final String signerAddress;

  @HiveField(6)
  final String merkleRoot;

  @HiveField(7)
  final DateTime anchoredAt;

  @HiveField(8)
  final String? location;

  @HiveField(9)
  final String? eventContext;

  @HiveField(10)
  final String? author;

  @HiveField(11)
  final String? license;

  @HiveField(12)
  final String imagePath;

  @HiveField(13)
  bool isOnChain;

  @HiveField(14)
  final String mediaType; // 'image' or 'video'

  AnchorRecord({
    required this.id,
    required this.sha512Hash,
    required this.pHash,
    required this.ipfsCID,
    required this.txHash,
    required this.signerAddress,
    required this.merkleRoot,
    required this.anchoredAt,
    this.location,
    this.eventContext,
    this.author,
    this.license,
    required this.imagePath,
    this.isOnChain = false,
    this.mediaType = 'image',
  });

  factory AnchorRecord.create({
    required String sha512Hash,
    required String pHash,
    required String ipfsCID,
    required String txHash,
    required String signerAddress,
    required String merkleRoot,
    String? location,
    String? eventContext,
    String? author,
    String? license,
    required String imagePath,
    bool isOnChain = false,
    String mediaType = 'image',
  }) {
    return AnchorRecord(
      id: const Uuid().v4(),
      sha512Hash: sha512Hash,
      pHash: pHash,
      ipfsCID: ipfsCID,
      txHash: txHash,
      signerAddress: signerAddress,
      merkleRoot: merkleRoot,
      anchoredAt: DateTime.now(),
      location: location,
      eventContext: eventContext,
      author: author,
      license: license,
      imagePath: imagePath,
      isOnChain: isOnChain,
      mediaType: mediaType,
    );
  }
}
