import 'package:flutter/material.dart';
import '../theme/app_colors.dart';

class MetricBar extends StatelessWidget {
  final String label;
  final String valueText;
  final Color valueColor;
  final double progress;
  final String note;

  const MetricBar({
    super.key,
    required this.label,
    required this.valueText,
    required this.valueColor,
    required this.progress,
    required this.note,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Text(
              label,
              style: const TextStyle(
                fontFamily: 'DM Sans',
                fontSize: 12,
                color: AppColors.textMuted,
              ),
            ),
            const Spacer(),
            Text(
              valueText,
              style: TextStyle(
                fontFamily: 'Space Mono',
                fontSize: 12,
                color: valueColor,
              ),
            ),
          ],
        ),
        const SizedBox(height: 3),
        Text(
          note,
          style: const TextStyle(
            fontFamily: 'Space Mono',
            fontSize: 10,
            color: AppColors.textHint,
          ),
        ),
        const SizedBox(height: 6),
        ClipRRect(
          borderRadius: BorderRadius.circular(2),
          child: TweenAnimationBuilder<double>(
            tween: Tween<double>(begin: 0, end: progress),
            duration: const Duration(milliseconds: 800),
            curve: Curves.easeInOut,
            builder: (context, value, child) {
              return LinearProgressIndicator(
                value: value,
                backgroundColor: AppColors.border,
                valueColor: AlwaysStoppedAnimation<Color>(valueColor),
                minHeight: 3,
              );
            },
          ),
        ),
        const SizedBox(height: 16),
      ],
    );
  }
}
