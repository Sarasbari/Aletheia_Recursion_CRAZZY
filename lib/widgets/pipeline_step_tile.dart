import 'package:flutter/material.dart';
import '../theme/app_colors.dart';

enum PipelineStatus { pending, active, done, error }

class PipelineStepTile extends StatelessWidget {
  final int stepNumber;
  final String title;
  final String value;
  final PipelineStatus status;

  const PipelineStepTile({
    super.key,
    required this.stepNumber,
    required this.title,
    required this.value,
    required this.status,
  });

  @override
  Widget build(BuildContext context) {
    final bool isDone = status == PipelineStatus.done;
    final bool isActive = status == PipelineStatus.active;
    final bool isPending = status == PipelineStatus.pending;

    return Container(
      padding: const EdgeInsets.symmetric(vertical: 14),
      decoration: const BoxDecoration(
        border: Border(bottom: BorderSide(color: AppColors.border, width: 1)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 24,
            child: isActive
                ? const Center(
                    child: SizedBox(
                      width: 12,
                      height: 12,
                      child: CircularProgressIndicator(
                        strokeWidth: 1.5,
                        color: AppColors.green,
                      ),
                    ),
                  )
                : Text(
                    '0$stepNumber',
                    style: TextStyle(
                      fontFamily: 'Space Mono',
                      fontSize: 10,
                      color: isDone ? AppColors.green : (isPending ? AppColors.textHint : AppColors.textMuted),
                    ),
                  ),
          ),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: TextStyle(
                    fontFamily: 'DM Sans',
                    fontSize: 13,
                    color: isDone ? AppColors.textMuted : (isActive ? AppColors.textPrimary : AppColors.textHint),
                  ),
                ),
                const SizedBox(width: 3),
                AnimatedDefaultTextStyle(
                  duration: const Duration(milliseconds: 300),
                  style: TextStyle(
                    fontFamily: 'Space Mono',
                    fontSize: 10,
                    color: isDone ? AppColors.green : (isActive ? AppColors.textMuted : AppColors.textHint),
                  ),
                  child: Text(value),
                ),
              ],
            ),
          ),
          if (isDone)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
              decoration: BoxDecoration(
                color: AppColors.greenDark,
                border: Border.all(color: AppColors.greenBorder),
                borderRadius: BorderRadius.circular(3),
              ),
              child: const Text(
                'DONE',
                style: TextStyle(
                  fontFamily: 'Space Mono',
                  fontSize: 9,
                  color: AppColors.green,
                ),
              ),
            )
          else if (isPending)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
              decoration: BoxDecoration(
                color: AppColors.surfaceAlt,
                border: Border.all(color: AppColors.border),
                borderRadius: BorderRadius.circular(3),
              ),
              child: const Text(
                'PENDING',
                style: TextStyle(
                  fontFamily: 'Space Mono',
                  fontSize: 9,
                  color: AppColors.textHint,
                ),
              ),
            ),
        ],
      ),
    );
  }
}
