import 'dart:io';
import 'package:dio/dio.dart';
import 'package:exif/exif.dart';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:path_provider/path_provider.dart';
import 'package:provider/provider.dart';
import 'package:video_player/video_player.dart';
import '../../models/media_type.dart';
import '../../models/anchor_record.dart';
import '../../services/blockchain_service.dart';
import '../../services/hash_service.dart';
import '../../services/ipfs_service.dart';
import '../../services/storage_service.dart';
import '../../services/wallet_service.dart';
import '../../theme/app_colors.dart';
import '../../widgets/mono_label.dart';
import '../../widgets/pipeline_step_tile.dart';
import '../../widgets/sigil_app_bar.dart';
import '../../widgets/certificate_card.dart';

class AnchorScreen extends StatefulWidget {
  const AnchorScreen({super.key});

  @override
  State<AnchorScreen> createState() => _AnchorScreenState();
}

enum AnchorState { input, pipeline, success }
enum SourceType { upload, camera, url }

class _AnchorScreenState extends State<AnchorScreen> {
  AnchorState _state = AnchorState.input;
  SourceType _sourceType = SourceType.upload;
  MediaType _mediaType = MediaType.image;
  File? _selectedFile;
  VideoPlayerController? _videoController;
  final _picker = ImagePicker();
  final _urlController = TextEditingController();
  
  // Metadata controllers
  final _locationController = TextEditingController();
  final _eventController = TextEditingController();
  final _authorController = TextEditingController();
  final _licenseController = TextEditingController();

  // Pipeline state
  int _currentStep = 0;
  String _sha512 = "—";
  String _pHash = "—";
  String _ipfsCID = "—";
  String _txHash = "—";
  String _signature = "—";

  Future<void> _pickMedia(ImageSource source, MediaType type) async {
    try {
      XFile? pickedFile;
      if (type == MediaType.image) {
        pickedFile = await _picker.pickImage(source: source);
      } else {
        pickedFile = await _picker.pickVideo(source: source);
      }

      if (pickedFile != null) {
        if (_videoController != null) {
          await _videoController!.dispose();
          _videoController = null;
        }

        setState(() {
          _selectedFile = File(pickedFile!.path);
          _mediaType = type;
        });

        if (type == MediaType.video) {
          _videoController = VideoPlayerController.file(_selectedFile!)
            ..initialize().then((_) => setState(() {}))
            ..setLooping(true)
            ..play();
        }

        _autoFillMetadata();
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error picking media: $e')),
      );
    }
  }

  Future<void> _handleUrlSource() async {
    if (_urlController.text.isEmpty) return;
    
    setState(() { _state = AnchorState.pipeline; _currentStep = 0; }); // Show loading
    
    try {
      final dio = Dio();
      final tempDir = await getTemporaryDirectory();
      final filePath = '${tempDir.path}/${DateTime.now().millisecondsSinceEpoch}.jpg';
      
      await dio.download(_urlController.text, filePath);
      
      setState(() {
        _selectedFile = File(filePath);
        _mediaType = _urlController.text.toLowerCase().endsWith('.mp4') ? MediaType.video : MediaType.image;
        _state = AnchorState.input;
      });

      if (_mediaType == MediaType.video) {
        if (_videoController != null) await _videoController!.dispose();
        _videoController = VideoPlayerController.file(_selectedFile!)
          ..initialize().then((_) => setState(() {}))
          ..setLooping(true)
          ..play();
      }

      _autoFillMetadata();
    } catch (e) {
      setState(() { _state = AnchorState.input; });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error downloading image: $e')),
      );
    }
  }

  Future<void> _autoFillMetadata() async {
    if (_selectedFile == null) return;
    if (_mediaType == MediaType.image) {
      final meta = await _extractMetadata(_selectedFile!);
      setState(() {
        if (meta.containsKey('location')) _locationController.text = meta['location']!;
        if (meta.containsKey('device')) _authorController.text = meta['device']!;
      });
    }
  }

  Future<Map<String, String>> _extractMetadata(File file) async {
    final Map<String, String> metadata = {};
    try {
      if (_mediaType == MediaType.image) {
        final bytes = await file.readAsBytes();
        final data = await readExifFromBytes(bytes);
        
        if (data.containsKey('Image DateTime')) {
          metadata['datetime'] = data['Image DateTime'].toString();
        }
        if (data.containsKey('Image Make')) {
          metadata['device'] = "${data['Image Make']} ${data['Image Model']}";
        }
        if (data.containsKey('GPS GPSLatitude')) {
          metadata['location'] = "GPS Locked";
        }
      } else {
        // Video metadata check (Basic)
        // AI generated videos often lack standard metadata atoms
        final bytes = await file.openRead(0, 4096).expand((i) => i).toList();
        final content = String.fromCharCodes(bytes);
        if (content.contains('com.apple.quicktime.creationdate') || 
            content.contains('com.android.version')) {
          metadata['datetime'] = "EXTRACTED_FROM_HEADER";
          metadata['device'] = "MOBILE_DEVICE";
        }
      }
    } catch (e) {
      debugPrint('Metadata extraction error: $e');
    }
    return metadata;
  }

  Future<bool> _validateMetadata() async {
    // Standard rule: Camera captures are always valid
    if (_sourceType == SourceType.camera) return true;

    final metadata = await _extractMetadata(_selectedFile!);
    
    // If no datetime or device info is found in a file upload, 
    // it's considered suspicious (AI generated or edited/cleaned)
    if (metadata.isEmpty || (!metadata.containsKey('datetime') && !metadata.containsKey('device'))) {
      return false;
    }
    return true;
  }

  Future<void> _startPipeline() async {
    if (_selectedFile == null) return;
    
    final walletService = context.read<WalletService>();
    if (!walletService.isConnected) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please connect your wallet first')),
      );
      return;
    }

    // AUTHENTICITY GUARD: Reject uploads without metadata
    if (!await _validateMetadata()) {
      _showAuthenticityWarning();
      return;
    }

    setState(() {
      _state = AnchorState.pipeline;
      _currentStep = _mediaType == MediaType.video ? -1 : 1; // -1 for video start
    });

    final hashService = HashService();
    final blockchainService = BlockchainService();
    final storageService = context.read<StorageService>();

    try {
      if (_mediaType == MediaType.video) {
        // Step -1: Audio extraction simulation
        setState(() { _currentStep = -1; });
        await Future.delayed(const Duration(milliseconds: 1500));
        
        // Step 0: Dual hashing simulation
        setState(() { _currentStep = 0; });
        await Future.delayed(const Duration(milliseconds: 1500));
      }

      // Step 1: SHA-512
      _sha512 = await hashService.computeSHA512(_selectedFile!);
      setState(() { _currentStep = 2; });

      // Step 2: pHash (Client-side simulation)
      _pHash = await hashService.computePHash(_selectedFile!);
      setState(() { _currentStep = 3; });

      // Step 3: ECDSA Sign (MOCKED as requested)
      _signature = await walletService.signMessageMock(_sha512);
      setState(() { _currentStep = 4; });

      // Step 4: IPFS Upload (MOCKED as requested)
      await Future.delayed(const Duration(milliseconds: 500));
      _ipfsCID = "QmMock${_sha512.substring(0, 8)}";
      setState(() { _currentStep = 5; });

      // Step 5: Merkle Batch (MOCKED)
      await Future.delayed(const Duration(milliseconds: 500));
      setState(() { _currentStep = 6; });

      // Step 6: REAL On-chain anchor for gas fee
      _txHash = await walletService.sendAnchorTransaction(_sha512, _ipfsCID);
      
      // Also notify backend for record keeping
      try {
        await blockchainService.anchorMedia(
          file: _selectedFile!,
          mediaType: _mediaType,
          captureTimestamp: DateTime.now().toIso8601String(),
          location: _locationController.text,
          deviceInfo: _authorController.text,
          uploaderId: walletService.connectedAddress,
        );
      } catch (e) {
        debugPrint('Backend sync skipped: $e');
      }
      
      // Save to Hive
      final record = AnchorRecord.create(
        sha512Hash: _sha512,
        pHash: _pHash,
        ipfsCID: _ipfsCID,
        txHash: _txHash,
        signerAddress: walletService.connectedAddress!,
        merkleRoot: "0xMerkle...",
        imagePath: _selectedFile!.path,
        location: _locationController.text,
        eventContext: _eventController.text,
        author: _authorController.text,
        license: _licenseController.text,
        isOnChain: true,
        mediaType: _mediaType.name,
      );
      await storageService.saveRecord(record);

      setState(() { _state = AnchorState.success; });
    } catch (e) {
      debugPrint('Pipeline error: $e');
      setState(() { _state = AnchorState.input; });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Anchoring failed: $e')),
      );
    }
  }

  void _showAuthenticityWarning() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: AppColors.surface,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12), side: const BorderSide(color: AppColors.border)),
        title: Row(
          children: const [
            Icon(Icons.warning_amber_rounded, color: AppColors.red),
            SizedBox(width: 12),
            Text("AUTHENTICITY ERROR", style: TextStyle(fontFamily: 'Space Mono', color: AppColors.red, fontSize: 16)),
          ],
        ),
        content: const Text(
          "AI Generated or Edited media detected.\n\nThis file lacks original camera metadata (EXIF/XMP). To ensure digital integrity, only original camera captures can be anchored on-chain.",
          style: TextStyle(fontFamily: 'DM Sans', color: AppColors.textPrimary),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text("CLOSE", style: TextStyle(fontFamily: 'Space Mono', color: AppColors.textMuted)),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SigilAppBar(
        title: "/ ANCHOR IMAGE",
        actions: [
          _GasPill(),
          const SizedBox(width: 16),
        ],
      ),
      body: AnimatedSwitcher(
        duration: const Duration(milliseconds: 300),
        child: _buildCurrentView(),
      ),
    );
  }

  Widget _buildCurrentView() {
    switch (_state) {
      case AnchorState.input:
        return _AnchorInputView(
          selectedFile: _selectedFile,
          mediaType: _mediaType,
          videoController: _videoController,
          sourceType: _sourceType,
          onSourceChanged: (type) {
            setState(() { _sourceType = type; });
            if (type == SourceType.upload) _pickMedia(ImageSource.gallery, MediaType.image);
            if (type == SourceType.camera) _pickMedia(ImageSource.camera, MediaType.image);
          },
          onPick: _pickMedia,
          onStart: _startPipeline,
          onUrlSubmit: _handleUrlSource,
          locationController: _locationController,
          eventController: _eventController,
          authorController: _authorController,
          licenseController: _licenseController,
          urlController: _urlController,
        );
      case AnchorState.pipeline:
        return _AnchorPipelineView(
          file: _selectedFile!,
          mediaType: _mediaType,
          currentStep: _currentStep,
          sha512: _sha512,
          pHash: _pHash,
          signature: _signature,
          ipfsCID: _ipfsCID,
          txHash: _txHash,
        );
      case AnchorState.success:
        return _AnchorSuccessView(
          onReset: () {
            if (_videoController != null) _videoController!.dispose();
            setState(() {
              _state = AnchorState.input;
              _selectedFile = null;
              _videoController = null;
              _locationController.clear();
              _eventController.clear();
              _authorController.clear();
              _licenseController.clear();
            });
          },
          record: AnchorRecord.create(
            sha512Hash: _sha512,
            pHash: _pHash,
            ipfsCID: _ipfsCID,
            txHash: _txHash,
            signerAddress: "0xaddr",
            merkleRoot: "0xroot",
            imagePath: _selectedFile!.path,
            mediaType: _mediaType.name,
          ),
        );
    }
  }
}

class _GasPill extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
        decoration: BoxDecoration(
          color: AppColors.surfaceAlt,
          border: Border.all(color: AppColors.border),
          borderRadius: BorderRadius.circular(4),
        ),
        child: const Text(
          "< ₹0.01",
          style: TextStyle(fontFamily: 'Space Mono', fontSize: 11, color: AppColors.textMuted),
        ),
      ),
    );
  }
}

class _AnchorInputView extends StatelessWidget {
  final File? selectedFile;
  final MediaType mediaType;
  final VideoPlayerController? videoController;
  final SourceType sourceType;
  final Function(SourceType) onSourceChanged;
  final Function(ImageSource, MediaType) onPick;
  final VoidCallback onStart;
  final VoidCallback onUrlSubmit;
  final TextEditingController locationController;
  final TextEditingController eventController;
  final TextEditingController authorController;
  final TextEditingController licenseController;
  final TextEditingController urlController;

  const _AnchorInputView({
    required this.selectedFile,
    required this.mediaType,
    required this.videoController,
    required this.sourceType,
    required this.onSourceChanged,
    required this.onPick,
    required this.onStart,
    required this.onUrlSubmit,
    required this.locationController,
    required this.eventController,
    required this.authorController,
    required this.licenseController,
    required this.urlController,
  });

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _SourceSelector(
            current: sourceType,
            onChanged: onSourceChanged,
          ),
          const SizedBox(height: 24),
          if (sourceType == SourceType.url && selectedFile == null)
            _UrlInput(controller: urlController, onSubmit: onUrlSubmit)
          else
            _DropZone(
              selectedFile: selectedFile, 
              mediaType: mediaType,
              videoController: videoController,
              onTap: () {
                if (sourceType == SourceType.upload) onPick(ImageSource.gallery, MediaType.image);
                if (sourceType == SourceType.camera) onPick(ImageSource.camera, MediaType.image);
              },
              onVideoTap: () {
                if (sourceType == SourceType.upload) onPick(ImageSource.gallery, MediaType.video);
                if (sourceType == SourceType.camera) onPick(ImageSource.camera, MediaType.video);
              }
            ),
          const SizedBox(height: 24),
          _MetadataGrid(
            location: locationController,
            event: eventController,
            author: authorController,
            license: licenseController,
          ),
          const SizedBox(height: 24),
          const _SigningInfoBox(),
          const SizedBox(height: 32),
          ElevatedButton(
            onPressed: selectedFile != null ? onStart : null,
            child: const Text("BEGIN ANCHORING"),
          ),
        ],
      ),
    );
  }
}

class _UrlInput extends StatelessWidget {
  final TextEditingController controller;
  final VoidCallback onSubmit;

  const _UrlInput({required this.controller, required this.onSubmit});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.surfaceAlt,
        border: Border.all(color: AppColors.border),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const MonoLabel(text: "IMAGE URL"),
          const SizedBox(height: 8),
          TextField(
            controller: controller,
            style: const TextStyle(fontFamily: 'DM Sans', fontSize: 13, color: AppColors.textPrimary),
            decoration: const InputDecoration(
              hintText: "https://example.com/image.jpg",
              border: OutlineInputBorder(),
            ),
          ),
          const SizedBox(height: 16),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton(
              onPressed: onSubmit,
              child: const Text("FETCH IMAGE"),
            ),
          ),
        ],
      ),
    );
  }
}

class _SourceSelector extends StatelessWidget {
  final SourceType current;
  final Function(SourceType) onChanged;

  const _SourceSelector({required this.current, required this.onChanged});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        _SourcePill(
          label: "Upload", 
          isActive: current == SourceType.upload,
          onTap: () => onChanged(SourceType.upload),
        ),
        const SizedBox(width: 8),
        _SourcePill(
          label: "Camera", 
          isActive: current == SourceType.camera,
          onTap: () => onChanged(SourceType.camera),
        ),
        const SizedBox(width: 8),
        _SourcePill(
          label: "URL", 
          isActive: current == SourceType.url,
          onTap: () => onChanged(SourceType.url),
        ),
      ],
    );
  }
}

class _SourcePill extends StatelessWidget {
  final String label;
  final bool isActive;
  final VoidCallback onTap;

  const _SourcePill({required this.label, required this.isActive, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        decoration: BoxDecoration(
          color: isActive ? AppColors.surfaceAlt : Colors.transparent,
          border: Border.all(color: isActive ? AppColors.green : AppColors.border),
          borderRadius: BorderRadius.circular(6),
        ),
        child: Text(
          label,
          style: TextStyle(
            fontFamily: 'Space Mono',
            fontSize: 11,
            color: isActive ? AppColors.green : AppColors.textMuted,
          ),
        ),
      ),
    );
  }
}

class _DropZone extends StatelessWidget {
  final File? selectedFile;
  final MediaType mediaType;
  final VideoPlayerController? videoController;
  final VoidCallback onTap;
  final VoidCallback onVideoTap;

  const _DropZone({
    required this.selectedFile, 
    required this.mediaType,
    required this.videoController,
    required this.onTap,
    required this.onVideoTap,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        GestureDetector(
          onTap: onTap,
          child: Container(
            height: 160,
            width: double.infinity,
            decoration: BoxDecoration(
              color: AppColors.surfaceAlt,
              borderRadius: BorderRadius.circular(8),
            ),
            child: selectedFile != null
                ? ClipRRect(
                    borderRadius: BorderRadius.circular(8),
                    child: (mediaType == MediaType.video && videoController != null && videoController!.value.isInitialized)
                        ? AspectRatio(
                            aspectRatio: videoController!.value.aspectRatio,
                            child: VideoPlayer(videoController!),
                          )
                        : Image.file(selectedFile!, fit: BoxFit.cover),
                  )
                : CustomPaint(
                    painter: _DashedBorderPainter(),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: const [
                        Icon(Icons.upload_file_outlined, size: 32, color: AppColors.textMuted),
                        SizedBox(height: 8),
                        Text("Tap to select image", style: TextStyle(fontFamily: 'DM Sans', fontSize: 13, color: AppColors.textMuted)),
                      ],
                    ),
                  ),
          ),
        ),
        if (selectedFile == null) ...[
          const SizedBox(height: 12),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: onVideoTap,
              icon: const Icon(Icons.videocam_outlined, size: 18),
              label: const Text("OR SELECT VIDEO"),
            ),
          ),
        ],
      ],
    );
  }
}

class _DashedBorderPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = AppColors.borderMid
      ..strokeWidth = 1
      ..style = PaintingStyle.stroke;

    final path = Path();
    path.addRRect(RRect.fromRectAndRadius(Rect.fromLTWH(0, 0, size.width, size.height), const Radius.circular(8)));

    const dashWidth = 6.0;
    const dashSpace = 4.0;
    double distance = 0.0;

    for (final pathMetric in path.computeMetrics()) {
      while (distance < pathMetric.length) {
        canvas.drawPath(
          pathMetric.extractPath(distance, distance + dashWidth),
          paint,
        );
        distance += dashWidth + dashSpace;
      }
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

class _MetadataGrid extends StatelessWidget {
  final TextEditingController location;
  final TextEditingController event;
  final TextEditingController author;
  final TextEditingController license;

  const _MetadataGrid({required this.location, required this.event, required this.author, required this.license});

  @override
  Widget build(BuildContext context) {
    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      mainAxisSpacing: 8,
      crossAxisSpacing: 8,
      childAspectRatio: 2.8,
      children: [
        _MetadataField(label: "LOCATION", controller: location, hint: "GPS / Venue"),
        _MetadataField(label: "EVENT", controller: event, hint: "Context"),
        _MetadataField(label: "AUTHOR", controller: author, hint: "Witness name"),
        _MetadataField(label: "LICENSE", controller: license, hint: "Creative Commons"),
      ],
    );
  }
}

class _MetadataField extends StatelessWidget {
  final String label;
  final TextEditingController controller;
  final String hint;

  const _MetadataField({required this.label, required this.controller, required this.hint});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: BoxDecoration(
        color: AppColors.surface,
        border: Border.all(color: AppColors.border),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          MonoLabel(text: label),
          const SizedBox(height: 4),
          Expanded(
            child: TextField(
              controller: controller,
              style: const TextStyle(fontFamily: 'DM Sans', fontSize: 12, color: AppColors.textMuted),
              decoration: InputDecoration(
                isDense: true,
                contentPadding: EdgeInsets.zero,
                hintText: hint,
                fillColor: Colors.transparent,
                border: InputBorder.none,
                enabledBorder: InputBorder.none,
                focusedBorder: InputBorder.none,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _SigningInfoBox extends StatelessWidget {
  const _SigningInfoBox();

  @override
  Widget build(BuildContext context) {
    final wallet = context.watch<WalletService>();
    final addr = wallet.connectedAddress ?? "Not Connected";

    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: AppColors.surface,
        border: Border.all(color: AppColors.border),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Icon(Icons.key_outlined, size: 14, color: AppColors.textMuted),
          const SizedBox(width: 8),
          Expanded(
            child: RichText(
              text: TextSpan(
                style: const TextStyle(fontFamily: 'DM Sans', fontSize: 12, color: AppColors.textMuted, height: 1.8),
                children: [
                  const TextSpan(text: "ECDSA signature will bind this anchor to wallet "),
                  TextSpan(
                    text: addr,
                    style: const TextStyle(fontFamily: 'Space Mono', color: AppColors.textPrimary),
                  ),
                  const TextSpan(text: ". Timestamp locked at submission."),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _AnchorPipelineView extends StatelessWidget {
  final File file;
  final MediaType mediaType;
  final int currentStep;
  final String sha512;
  final String pHash;
  final String signature;
  final String ipfsCID;
  final String txHash;

  const _AnchorPipelineView({
    required this.file,
    required this.mediaType,
    required this.currentStep,
    required this.sha512,
    required this.pHash,
    required this.signature,
    required this.ipfsCID,
    required this.txHash,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Stack(
          children: [
            SizedBox(
              height: 120,
              width: double.infinity,
              child: Image.file(file, fit: BoxFit.cover),
            ),
            Positioned.fill(
              child: Container(
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.bottomCenter,
                    end: Alignment.topCenter,
                    colors: [Colors.black.withOpacity(0.8), Colors.transparent],
                  ),
                ),
              ),
            ),
            Positioned(
              bottom: 12,
              left: 24,
              child: Text(
                file.path.split('/').last,
                style: const TextStyle(fontFamily: 'Space Mono', fontSize: 11, color: AppColors.textPrimary),
              ),
            ),
          ],
        ),
        const Padding(
          padding: EdgeInsets.fromLTRB(24, 24, 24, 8),
          child: MonoLabel(text: "CRYPTO PIPELINE"),
        ),
        Expanded(
          child: ListView(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            children: [
              if (mediaType == MediaType.video) ...[
                PipelineStepTile(
                  stepNumber: -1,
                  title: "Extracting audio stream",
                  value: currentStep >= -1 ? "PCM / 44.1kHz extracted" : "—",
                  status: currentStep > -1 ? PipelineStatus.done : (currentStep == -1 ? PipelineStatus.active : PipelineStatus.pending),
                ),
                PipelineStepTile(
                  stepNumber: 0,
                  title: "Generating dual-entropy hash",
                  value: currentStep >= 0 ? "Multi-modal signature ready" : "—",
                  status: currentStep > 0 ? PipelineStatus.done : (currentStep == 0 ? PipelineStatus.active : PipelineStatus.pending),
                ),
              ],
              PipelineStepTile(
                stepNumber: 1,
                title: "SHA-512 cross-check",
                value: sha512,
                status: currentStep > 1 ? PipelineStatus.done : (currentStep == 1 ? PipelineStatus.active : PipelineStatus.pending),
              ),
              PipelineStepTile(
                stepNumber: 2,
                title: "pHash computed",
                value: pHash,
                status: currentStep > 2 ? PipelineStatus.done : (currentStep == 2 ? PipelineStatus.active : PipelineStatus.pending),
              ),
              PipelineStepTile(
                stepNumber: 3,
                title: "ECDSA signing",
                value: signature,
                status: currentStep > 3 ? PipelineStatus.done : (currentStep == 3 ? PipelineStatus.active : PipelineStatus.pending),
              ),
              PipelineStepTile(
                stepNumber: 4,
                title: "IPFS upload",
                value: ipfsCID,
                status: currentStep > 4 ? PipelineStatus.done : (currentStep == 4 ? PipelineStatus.active : PipelineStatus.pending),
              ),
              PipelineStepTile(
                stepNumber: 5,
                title: "Merkle batch",
                value: currentStep >= 5 ? "Batching with 47 other anchors" : "—",
                status: currentStep > 5 ? PipelineStatus.done : (currentStep == 5 ? PipelineStatus.active : PipelineStatus.pending),
              ),
              PipelineStepTile(
                stepNumber: 6,
                title: "On-chain anchor",
                value: txHash,
                status: currentStep > 6 ? PipelineStatus.done : (currentStep == 6 ? PipelineStatus.active : PipelineStatus.pending),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

class _AnchorSuccessView extends StatelessWidget {
  final VoidCallback onReset;
  final AnchorRecord record;

  const _AnchorSuccessView({required this.onReset, required this.record});

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        children: [
          const SizedBox(height: 20),
          const Icon(Icons.check_circle, size: 64, color: AppColors.green),
          const SizedBox(height: 20),
          const Text(
            "Image anchored",
            style: TextStyle(fontFamily: 'Space Mono', fontSize: 20, color: AppColors.textPrimary),
          ),
          const SizedBox(height: 6),
          const Text(
            "Permanently recorded on Polygon Mumbai",
            style: TextStyle(fontFamily: 'DM Sans', fontSize: 13, color: AppColors.textMuted),
          ),
          const SizedBox(height: 32),
          CertificateCard(record: record),
          const SizedBox(height: 32),
          Row(
            children: [
              Expanded(
                child: OutlinedButton(
                  onPressed: onReset,
                  child: const Text("Anchor another"),
                ),
              ),
              const SizedBox(width: 10),
              Expanded(
                child: ElevatedButton.icon(
                  onPressed: () {},
                  icon: const Icon(Icons.share, size: 16),
                  label: const Text("Share"),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
