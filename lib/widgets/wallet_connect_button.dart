import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../services/wallet_service.dart';
import '../theme/app_colors.dart';

class WalletConnectButton extends StatefulWidget {
  const WalletConnectButton({super.key});

  @override
  State<WalletConnectButton> createState() => _WalletConnectButtonState();
}

class _WalletConnectButtonState extends State<WalletConnectButton> with SingleTickerProviderStateMixin {
  late AnimationController _pulseController;

  @override
  void initState() {
    super.initState();
    _pulseController = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 2),
    )..repeat(reverse: true);
  }

  @override
  void dispose() {
    _pulseController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final walletService = context.watch<WalletService>();

    if (!walletService.isConnected) {
      return OutlinedButton(
        onPressed: () => walletService.connectWallet(context),
        style: OutlinedButton.styleFrom(
          side: const BorderSide(color: AppColors.borderMid),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          minimumSize: Size.zero,
          foregroundColor: AppColors.textMuted,
        ),
        child: const Text(
          'CONNECT WALLET',
          style: TextStyle(fontFamily: 'Space Mono', fontSize: 11),
        ),
      );
    }

    final address = walletService.connectedAddress ?? '';
    final shortAddress = address.length >= 10
        ? '${address.substring(0, 6)}...${address.substring(address.length - 4)}'
        : address;

    return GestureDetector(
      onTap: () => walletService.disconnect(),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
        decoration: BoxDecoration(
          color: AppColors.surfaceAlt,
          border: Border.all(color: AppColors.border),
          borderRadius: BorderRadius.circular(4),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            FadeTransition(
              opacity: Tween<double>(begin: 0.3, end: 1.0).animate(_pulseController),
              child: Container(
                width: 6,
                height: 6,
                decoration: const BoxDecoration(
                  color: AppColors.green,
                  shape: BoxShape.circle,
                ),
              ),
            ),
            const SizedBox(width: 8),
            Text(
              shortAddress,
              style: const TextStyle(
                color: AppColors.textMuted,
                fontFamily: 'Space Mono',
                fontSize: 11,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
