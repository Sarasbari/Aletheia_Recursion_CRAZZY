import 'dart:io';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:intl/intl.dart';
import '../../models/anchor_record.dart';
import '../../models/media_type.dart';
import '../../services/storage_service.dart';
import '../../theme/app_colors.dart';
import '../../widgets/sigil_app_bar.dart';
import 'certificate_screen.dart';

class HistoryScreen extends StatelessWidget {
  const HistoryScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final records = context.watch<StorageService>().getAllRecords();

    return Scaffold(
      appBar: SigilAppBar(
        title: "/ MY ANCHORS",
        actions: [
          IconButton(
            icon: const Icon(Icons.filter_list, color: AppColors.textMuted),
            onPressed: () {},
          ),
        ],
      ),
      body: records.isEmpty 
        ? const _EmptyHistoryView()
        : ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: records.length,
            itemBuilder: (context, index) => _AnchorRecordTile(record: records[index]),
          ),
    );
  }
}

class _EmptyHistoryView extends StatelessWidget {
  const _EmptyHistoryView();

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: const [
          Icon(Icons.history_outlined, size: 48, color: AppColors.textHint),
          SizedBox(height: 16),
          Text(
            "No anchors yet",
            style: TextStyle(fontFamily: 'Space Mono', fontSize: 16, color: AppColors.textMuted),
          ),
          Text(
            "Images you anchor will appear here",
            style: TextStyle(fontFamily: 'DM Sans', fontSize: 13, color: AppColors.textHint),
          ),
        ],
      ),
    );
  }
}

class _AnchorRecordTile extends StatelessWidget {
  final AnchorRecord record;

  const _AnchorRecordTile({required this.record});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () => Navigator.push(
        context,
        MaterialPageRoute(builder: (_) => CertificateScreen(record: record)),
      ),
      child: Container(
        margin: const EdgeInsets.only(bottom: 8),
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: AppColors.surface,
          border: Border.all(color: AppColors.border),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: SizedBox(
                width: 56,
                height: 56,
                child: (record.mediaType == 'video')
                  ? Container(
                      color: AppColors.surfaceAlt,
                      child: const Icon(Icons.videocam, color: AppColors.green),
                    )
                  : Image.file(
                      File(record.imagePath),
                      fit: BoxFit.cover,
                      errorBuilder: (context, error, stackTrace) => Container(
                        color: AppColors.surfaceAlt,
                        child: const Icon(Icons.image_outlined, color: AppColors.textHint),
                      ),
                    ),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          record.ipfsCID,
                          style: const TextStyle(fontFamily: 'Space Mono', fontSize: 11, color: AppColors.green),
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        DateFormat('MMM dd, yyyy').format(record.anchoredAt),
                        style: const TextStyle(fontFamily: 'DM Sans', fontSize: 11, color: AppColors.textMuted),
                      ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Text(
                    record.eventContext ?? 'No context',
                    style: const TextStyle(fontFamily: 'DM Sans', fontSize: 12, color: AppColors.textPrimary),
                  ),
                  const SizedBox(height: 6),
                  Row(
                    children: [
                      _StatusPill(isOnChain: record.isOnChain),
                      const SizedBox(width: 6),
                      Expanded(
                        child: Text(
                          record.txHash,
                          style: const TextStyle(fontFamily: 'Space Mono', fontSize: 10, color: AppColors.textHint),
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const Icon(Icons.arrow_forward_ios, size: 12, color: AppColors.textHint),
          ],
        ),
      ),
    );
  }
}

class _StatusPill extends StatelessWidget {
  final bool isOnChain;
  const _StatusPill({required this.isOnChain});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: isOnChain ? AppColors.greenDark : Colors.amber.withOpacity(0.1),
        border: Border.all(color: isOnChain ? AppColors.greenBorder : Colors.amber.withOpacity(0.5)),
        borderRadius: BorderRadius.circular(3),
      ),
      child: Text(
        isOnChain ? "ON-CHAIN" : "PENDING",
        style: TextStyle(
          fontFamily: 'Space Mono',
          fontSize: 9,
          color: isOnChain ? AppColors.green : Colors.amber,
        ),
      ),
    );
  }
}
