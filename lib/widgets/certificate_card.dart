import 'dart:convert';
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:qr_flutter/qr_flutter.dart';
import 'package:url_launcher/url_launcher.dart';
import '../models/anchor_record.dart';
import '../theme/app_colors.dart';
import 'mono_label.dart';

class CertificateCard extends StatelessWidget {
  final AnchorRecord record;

  const CertificateCard({super.key, required this.record});

  void _viewOnExplorer() async {
    final url = Uri.parse('https://mumbai.polygonscan.com/tx/${record.txHash}');
    if (await canLaunchUrl(url)) {
      await launchUrl(url, mode: LaunchMode.externalApplication);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      decoration: BoxDecoration(
        color: AppColors.surface,
        border: Border.all(color: AppColors.border, width: 1),
        borderRadius: BorderRadius.circular(12),
      ),
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              RichText(
                text: const TextSpan(
                  children: [
                    TextSpan(
                      text: "SIGIL",
                      style: TextStyle(
                        fontFamily: 'Space Mono',
                        fontSize: 13,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textPrimary,
                      ),
                    ),
                    TextSpan(
                      text: ".",
                      style: TextStyle(
                        fontFamily: 'Space Mono',
                        fontSize: 13,
                        fontWeight: FontWeight.bold,
                        color: AppColors.green,
                      ),
                    ),
                  ],
                ),
              ),
              const Spacer(),
              GestureDetector(
                onTap: _viewOnExplorer,
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: AppColors.greenDark,
                    border: Border.all(color: AppColors.greenBorder, width: 1),
                    borderRadius: BorderRadius.circular(3),
                  ),
                  child: Row(
                    children: const [
                      Text(
                        "ON-CHAIN ",
                        style: TextStyle(fontFamily: 'Space Mono', fontSize: 10, color: AppColors.green),
                      ),
                      Icon(Icons.open_in_new, size: 10, color: AppColors.green),
                    ],
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          ClipRRect(
            borderRadius: BorderRadius.circular(8),
            child: Container(
              height: 180,
              width: double.infinity,
              decoration: const BoxDecoration(color: AppColors.surfaceAlt),
              child: record.mediaType == 'video'
                  ? const Center(
                      child: Icon(Icons.videocam, size: 64, color: AppColors.green),
                    )
                  : Image.file(
                      File(record.imagePath),
                      fit: BoxFit.cover,
                      errorBuilder: (context, error, stackTrace) => Center(
                        child: Text(
                          record.sha512Hash.substring(0, 32),
                          style: const TextStyle(fontFamily: 'Space Mono', fontSize: 9, color: AppColors.textHint),
                          textAlign: TextAlign.center,
                        ),
                      ),
                    ),
            ),
          ),
          const SizedBox(height: 16),
          const Divider(),
          const SizedBox(height: 12),
          Wrap(
            spacing: 12,
            runSpacing: 8,
            children: [
              _DetailChip(label: "SHA-512", value: "${record.sha512Hash.substring(0, 8)}...${record.sha512Hash.substring(record.sha512Hash.length-4)}"),
              _DetailChip(label: "IPFS CID", value: "${record.ipfsCID.length > 12 ? record.ipfsCID.substring(0, 12) : record.ipfsCID}..."),
              _DetailChip(label: "TX HASH", value: "${record.txHash.length > 12 ? record.txHash.substring(0, 12) : record.txHash}..."),
              _DetailChip(label: "BLOCK", value: record.id.substring(0, 8)),
              _DetailChip(label: "ANCHORED", value: record.anchoredAt.toIso8601String().split('T').first),
              _DetailChip(label: "DEVICE", value: record.id.substring(record.id.length-4).toUpperCase()),
            ],
          ),
          const SizedBox(height: 16),
          const Divider(),
          const SizedBox(height: 12),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    "Scan to verify",
                    style: TextStyle(fontFamily: 'Space Mono', fontSize: 10, color: AppColors.textMuted),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    "Sigil · ${record.location ?? 'Global Network'}",
                    style: const TextStyle(fontFamily: 'DM Sans', fontSize: 11, color: AppColors.textHint),
                  ),
                ],
              ),
              QrImageView(
                data: jsonEncode({
                  'cid': record.ipfsCID,
                  'hash': record.sha512Hash,
                  'tx': record.txHash,
                }),
                size: 72,
                gapless: false,
                foregroundColor: AppColors.green,
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _DetailChip extends StatelessWidget {
  final String label;
  final String value;

  const _DetailChip({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label.toUpperCase(),
          style: const TextStyle(fontFamily: 'Space Mono', fontSize: 9, color: AppColors.textHint),
        ),
        const SizedBox(height: 2),
        Text(
          value,
          style: const TextStyle(fontFamily: 'Space Mono', fontSize: 11, color: AppColors.green),
        ),
      ],
    );
  }
}
