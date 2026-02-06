import 'package:echo_chat_app_frontend/data/models/user_model.dart';
import 'package:echo_chat_app_frontend/data/services/auth_service.dart';

class AuthRepository {
  final AuthService authService;

  AuthRepository({required this.authService});

  // Stream buatan sendiri untuk menggabungkan authStateChanges dari firebase dan syncUser
  // ketika ada perubahan status auth firebase, maka akan di sync ke backend
  // jika sync gagal, maka akan di logout
  Stream<User?> get authStateChanges {
    return authService.authStateChanges.asyncMap((firebaseUser) async {
      if (firebaseUser == null) return null;

      try {
        // Setiap kali status Firebase berubah (login),
        // kita otomatis ambil data lengkap dari Google
        return await authService.syncUser();
      } catch (e) {
        // Jika sync gagal, logout
        await authService.signOut();
        return null;
      }
    });
  }

  Future<void> signInWithGoogle() async {
    await authService.signInWithGoogle();
  }

  Future<void> signInWithEmailAndPassword(String email, String password) async {
    await authService.signInWithEmailAndPassword(email, password);
  }

  Future<void> signUpWithEmailAndPassword(String email, String password) async {
    await authService.signUpWithEmailAndPassword(email, password);
  }

  Future<void> signOut() async {
    await authService.signOut();
  }
}
