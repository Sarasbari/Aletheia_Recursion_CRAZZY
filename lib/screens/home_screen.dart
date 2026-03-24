import 'package:flutter/material.dart';
import '../theme/app_colors.dart';
import '../widgets/wallet_connect_button.dart';
import '../widgets/mono_label.dart';
import '../widgets/hash_ticker.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Top Bar
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 20),
                child: Row(
                  children: [
                    RichText(
                      text: const TextSpan(
                        children: [
                          TextSpan(
                            text: "SIGIL",
                            style: TextStyle(
                              fontFamily: 'Space Mono',
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                              color: AppColors.textPrimary,
                            ),
                          ),
                          TextSpan(
                            text: ".",
                            style: TextStyle(
                              fontFamily: 'Space Mono',
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                              color: AppColors.green,
                            ),
                          ),
                        ],
                      ),
                    ),
                    const Spacer(),
                    const WalletConnectButton(),
                  ],
                ),
              ),

              // Hero Content
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 32),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const MonoLabel(text: "CRYPTOGRAPHIC MEDIA TRUTH"),
                    const SizedBox(height: 12),
                    const HashTicker(),
                    const SizedBox(height: 24),
                    const Text(
                      "Truth has a\nfingerprint.",
                      style: TextStyle(
                        fontFamily: 'Space Mono',
                        fontSize: 26,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textPrimary,
                        height: 1.15,
                      ),
                    ),
                    const SizedBox(height: 12),
                    const Text(
                      "Anchor any media permanently.\nVerify any media instantly.",
                      style: TextStyle(
                        fontFamily: 'DM Sans',
                        fontSize: 14,
                        color: AppColors.textMuted,
                        height: 1.7,
                      ),
                    ),
                  ],
                ),
              ),

              // Action Cards
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24),
                child: Column(
                  children: [
                    _ActionCard(
                      tag: "01 — ANCHOR",
                      title: "Anchor your media",
                      subtitle: "Register images/videos on Polygon",
                      statsLine: "SHA-512 · ECDSA · IPFS · Polygon",
                      onTap: () {
                        // In a real app, this would use a provider to switch the shell index
                      },
                    ),
                    const SizedBox(height: 12),
                    _ActionCard(
                      tag: "02 — VERIFY",
                      title: "Verify media",
                      subtitle: "Check authenticity on-chain",
                      statsLine: "pHash · SSIM · Merkle · ZK",
                      onTap: () {},
                    ),
                  ],
                ),
              ),

              const SizedBox(height: 48),
              
              // Stats strip
              Container(
                decoration: const BoxDecoration(
                  border: Border(top: BorderSide(color: AppColors.border, width: 1)),
                ),
                padding: const EdgeInsets.symmetric(vertical: 24, horizontal: 24),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    _MetricColumn(value: "SHA-512", label: "Hash algo"),
                    _MetricColumn(value: "Polygon", label: "Blockchain"),
                    _MetricColumn(value: "< ₹0.01", label: "Gas per anchor"),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _ActionCard extends StatelessWidget {
  final String tag;
  final String title;
  final String subtitle;
  final String statsLine;
  final VoidCallback onTap;

  const _ActionCard({
    required this.tag,
    required this.title,
    required this.subtitle,
    required this.statsLine,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(12),
      child: Container(
        padding: const EdgeInsets.all(20),
        decoration: BoxDecoration(
          color: AppColors.surface,
          border: Border.all(color: AppColors.border, width: 1),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  MonoLabel(text: tag),
                  const SizedBox(height: 10),
                  Text(
                    title,
                    style: const TextStyle(
                      fontFamily: 'DM Sans',
                      fontSize: 16,
                      fontWeight: FontWeight.w500,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  Text(
                    subtitle,
                    style: const TextStyle(
                      fontFamily: 'DM Sans',
                      fontSize: 12,
                      color: AppColors.textMuted,
                    ),
                  ),
                  const SizedBox(height: 12),
                  Text(
                    statsLine,
                    style: const TextStyle(
                      fontFamily: 'Space Mono',
                      fontSize: 10,
                      color: AppColors.textHint,
                      letterSpacing: 0.06,
                    ),
                  ),
                ],
              ),
            ),
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                border: Border.all(color: AppColors.borderMid),
                borderRadius: BorderRadius.circular(6),
              ),
              child: const Icon(Icons.arrow_forward, color: AppColors.green, size: 18),
            ),
          ],
        ),
      ),
    );
  }
}

class _MetricColumn extends StatelessWidget {
  final String value;
  final String label;

  const _MetricColumn({required this.value, required this.label});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          value,
          style: const TextStyle(fontFamily: 'Space Mono', fontSize: 13, color: AppColors.textPrimary),
        ),
        const SizedBox(height: 4),
        Text(
          label.toUpperCase(),
          style: const TextStyle(fontFamily: 'Space Mono', fontSize: 10, color: Color(0xFF555555)),
        ),
      ],
    );
  }
}
