import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import '../theme/app_colors.dart';

class VerdictBanner extends StatelessWidget {
  final bool isAuthentic;
  final String blockNumber;
  final String date;
  final String statusText;
  final int trustScore;
  final String verdict;
  final String? txHash;

  const VerdictBanner({
    super.key,
    required this.isAuthentic,
    required this.blockNumber,
    required this.date,
    required this.statusText,
    required this.trustScore,
    required this.verdict,
    this.txHash,
  });

  void _viewOnExplorer() async {
    if (txHash == null || txHash!.isEmpty) return;
    final url = Uri.parse('https://mumbai.polygonscan.com/tx/$txHash');
    if (await canLaunchUrl(url)) {
      await launchUrl(url, mode: LaunchMode.externalApplication);
    }
  }

  @override
  Widget build(BuildContext context) {
    final bgColor = isAuthentic ? AppColors.greenDark : AppColors.redDark;
    final borderColor = isAuthentic ? AppColors.greenBorder : AppColors.redBorder;
    final accentColor = isAuthentic ? AppColors.green : AppColors.red;
    final icon = isAuthentic ? Icons.check_circle : Icons.warning_amber_rounded;
    
    // Web3 Accents
    final Color scoreColor = trustScore >= 80 ? AppColors.green : (trustScore >= 50 ? Colors.orange : AppColors.red);

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
      decoration: BoxDecoration(
        color: bgColor,
        border: Border(bottom: BorderSide(color: borderColor, width: 1)),
      ),
      child: Column(
        children: [
          Row(
            children: [
              Icon(icon, color: accentColor, size: 28),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      verdict.toUpperCase(),
                      style: TextStyle(
                        fontFamily: 'Space Mono',
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                        color: accentColor,
                        letterSpacing: 1.2,
                      ),
                    ),
                    Text(
                      statusText,
                      style: TextStyle(
                        fontFamily: 'DM Sans',
                        fontSize: 11,
                        color: accentColor.withOpacity(0.6),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 12),
              _TrustScoreIndicator(score: trustScore, color: scoreColor),
            ],
          ),
          const SizedBox(height: 12),
          const Divider(color: Colors.white10),
          Padding(
            padding: const EdgeInsets.only(top: 8),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                _MetaItem(label: "BLOCK", value: "#$blockNumber"),
                _MetaItem(label: "NETWORK", value: "ALETHEIA-P2P"),
                if (txHash != null && txHash!.length > 10)
                  GestureDetector(
                    onTap: _viewOnExplorer,
                    child: Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: Colors.white12,
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Row(
                        children: const [
                          Icon(Icons.open_in_new, size: 10, color: AppColors.green),
                          SizedBox(width: 4),
                          Text(
                            "EXPLORER",
                            style: TextStyle(fontFamily: 'Space Mono', fontSize: 8, color: AppColors.green),
                          ),
                        ],
                      ),
                    ),
                  )
                else
                  _MetaItem(label: "TIMESTAMP", value: date),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _TrustScoreIndicator extends StatelessWidget {
  final int score;
  final Color color;

  const _TrustScoreIndicator({required this.score, required this.color});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        border: Border.all(color: color.withOpacity(0.3), width: 2),
      ),
      child: Column(
        children: [
          Text(
            "$score",
            style: TextStyle(
              fontFamily: 'Space Mono',
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: color,
            ),
          ),
          Text(
            "TRUST",
            style: TextStyle(
              fontFamily: 'Space Mono',
              fontSize: 7,
              color: color.withOpacity(0.7),
            ),
          ),
        ],
      ),
    );
  }
}

class _MetaItem extends StatelessWidget {
  final String label;
  final String value;

  const _MetaItem({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(fontFamily: 'Space Mono', fontSize: 8, color: AppColors.textHint),
        ),
        Text(
          value,
          style: const TextStyle(fontFamily: 'Space Mono', fontSize: 10, color: AppColors.textMuted),
        ),
      ],
    );
  }
}
