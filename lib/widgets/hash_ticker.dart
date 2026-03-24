import 'package:flutter/material.dart';
import '../theme/app_colors.dart';

class HashTicker extends StatefulWidget {
  const HashTicker({super.key});

  @override
  State<HashTicker> createState() => _HashTickerState();
}

class _HashTickerState extends State<HashTicker> with SingleTickerProviderStateMixin {
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(seconds: 10),
    )..repeat();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ClipRect(
      child: AnimatedBuilder(
        animation: _controller,
        builder: (context, child) {
          return SizedBox(
            height: 20,
            child: OverflowBox(
              minWidth: 0.0,
              maxWidth: double.infinity,
              minHeight: 0.0,
              maxHeight: 20,
              alignment: Alignment.centerLeft,
              child: Transform.translate(
                offset: Offset(-200 * _controller.value, 0),
                child: Row(
                  children: List.generate(10, (index) => const Padding(
                    padding: EdgeInsets.only(right: 40),
                    child: Text(
                      '0x7f3a9c2b1e4d8f6a0c5b2e9d3f7a1c4e9f3a2b1e4d8f6a0c5b2e9d3f7a1c4e',
                      style: TextStyle(
                        fontFamily: 'Space Mono',
                        fontSize: 10,
                        color: AppColors.textHint,
                        overflow: TextOverflow.visible,
                      ),
                    ),
                  )),
                ),
              ),
            ),
          );
        },
      ),
    );
  }
}
