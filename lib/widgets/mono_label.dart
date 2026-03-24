import 'package:flutter/material.dart';
import '../theme/app_colors.dart';

class MonoLabel extends StatelessWidget {
  final String text;
  final Color color;

  const MonoLabel({
    super.key,
    required this.text,
    this.color = AppColors.green,
  });

  @override
  Widget build(BuildContext context) {
    return Text(
      text.toUpperCase(),
      style: TextStyle(
        fontFamily: 'Space Mono',
        fontSize: 10,
        color: color,
        letterSpacing: 0.12,
      ),
    );
  }
}
