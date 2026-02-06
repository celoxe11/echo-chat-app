import 'dart:io';
import 'package:flutter/foundation.dart';

class PlatformConfig {
  /// Base URL defined via --dart-define=BASE_URL=https://your-api.com
  /// If not provided, it falls back to the hardcoded defaults below.
  static const String _baseUrlFromEnv = String.fromEnvironment('BASE_URL');

  // Hardcoded Production URL (fallback)
  static const String _prodBaseUrl = 'https://api.echo-chat.com';

  // Development URLs
  // Android Emulator: 10.0.2.2 is mapped to your computer's localhost
  static const String _devBaseUrlAndroid = 'http://10.0.2.2:8080';
  // iOS Simulator, Web, Desktop: localhost
  static const String _devBaseUrlLocal = 'http://localhost:8080';

  /// Check if the app is running in production (release mode)
  static bool get isProduction => kReleaseMode;

  /// Get the base URL for the backend API based on environment and platform
  static String getBaseUrl() {
    // 1. If a URL is provided via --dart-define, use it (useful for CI/CD or staging)
    if (_baseUrlFromEnv.isNotEmpty) {
      return _baseUrlFromEnv;
    }

    // 2. If in production, use the production fallback
    if (isProduction) {
      return _prodBaseUrl;
    }

    // 3. Handle Development environment
    // For Web, dart:io is not available, so we check kIsWeb first
    if (kIsWeb) {
      return _devBaseUrlLocal;
    }

    // For Mobile/Desktop, we check the specific platform
    try {
      if (Platform.isAndroid) {
        return _devBaseUrlAndroid;
      }
    } catch (e) {
      // Fallback if Platform call fails (though kIsWeb coverage should prevent this)
    }

    return _devBaseUrlLocal;
  }

  /// Helper to get the WebSocket URL if needed
  static String getWebSocketUrl() {
    final baseUrl = getBaseUrl();
    if (baseUrl.startsWith('https')) {
      return baseUrl.replaceFirst('https', 'wss');
    }
    return baseUrl.replaceFirst('http', 'ws');
  }
}
