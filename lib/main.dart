import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:provider/provider.dart';
import 'models/anchor_record.dart';
import 'services/wallet_service.dart';
import 'services/storage_service.dart';
import 'theme/app_theme.dart';
import 'screens/main_shell.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Initialize Storage
  final storageService = StorageService();
  await storageService.init();
  
  // Initialize Wallet
  final walletService = WalletService();
  await walletService.initialize();
  
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider.value(value: walletService),
        Provider.value(value: storageService),
      ],
      child: const SigilApp(),
    ),
  );
}

class SigilApp extends StatelessWidget {
  const SigilApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Sigil',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.darkTheme,
      home: const MainShell(),
    );
  }
}
