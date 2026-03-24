import 'package:flutter/material.dart';
import '../theme/app_colors.dart';

class SigilAppBar extends StatelessWidget implements PreferredSizeWidget {
  final String title;
  final List<Widget>? actions;
  final Widget? leading;

  const SigilAppBar({
    super.key,
    required this.title,
    this.actions,
    this.leading,
  });

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text(
        title,
        style: const TextStyle(
          color: AppColors.textMuted,
          fontFamily: 'Space Mono',
          fontSize: 13,
        ),
      ),
      leading: leading ?? (Navigator.canPop(context) 
        ? IconButton(
            icon: const Icon(Icons.arrow_back, color: AppColors.textMuted, size: 20),
            onPressed: () => Navigator.pop(context),
          )
        : null),
      actions: actions,
      backgroundColor: AppColors.background,
      elevation: 0,
      centerTitle: true,
      shape: const Border(
        bottom: BorderSide(color: AppColors.border, width: 1),
      ),
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);
}
