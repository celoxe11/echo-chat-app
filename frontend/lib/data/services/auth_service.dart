import 'package:echo_chat_app_frontend/core/api/dio_client.dart';
import 'package:firebase_auth/firebase_auth.dart' as auth;
import 'package:google_sign_in/google_sign_in.dart';
import 'package:echo_chat_app_frontend/data/models/user_model.dart';
import 'dart:developer' as dev;

class AuthException implements Exception {
  final String message;
  AuthException(this.message);

  @override
  String toString() => message;
}

class AuthService {
  final auth.FirebaseAuth _auth = auth.FirebaseAuth.instance;
  final GoogleSignIn _googleSignIn = GoogleSignIn();
  final DioClient dioClient = DioClient.instance;

  // Get current user (user firebase)
  auth.User? get currentUser => _auth.currentUser;
  
  // Stream of auth state changes
  Stream<auth.User?> get authStateChanges => _auth.authStateChanges();

  // Sync user firebase ke backend, mereturn User model 
  Future<User> syncUser() async {
    final user = _auth.currentUser;
    if (user == null) throw AuthException('User not found');

    try {
      final response = await dioClient.post(
        '/auth/sync',
        data: {
          'uid': user.uid,
          'email': user.email,
          'displayName': user.displayName,
          'photoUrl': user.photoURL,
        },
      );
      return User.fromJson(response.data);
    } catch (e, stackTrace) {
      dev.log(
        'Error syncing user to backend',
        name: 'AuthService',
        error: e,
        stackTrace: stackTrace,
      );
      rethrow;
    }
  }

  Future<void> signInWithGoogle() async {
    try {
      // Sign out from previous sessions
      try {
        await _googleSignIn.signOut();
        await _auth.signOut();
      } catch (_) {}

      final googleUser = await _googleSignIn.signIn();
      if (googleUser == null) return; // user cancelled

      final googleAuth = await googleUser.authentication;
      final credential = auth.GoogleAuthProvider.credential(
        accessToken: googleAuth.accessToken,
        idToken: googleAuth.idToken,
      );

      await _auth.signInWithCredential(credential);
      await syncUser(); // trigger syncUser to backend
    } on auth.FirebaseAuthException catch (e, stackTrace) {
      final message = _handleFirebaseAuthError(e);
      dev.log(
        'Error signing in with Google: $message',
        name: 'AuthService',
        error: e,
        stackTrace: stackTrace,
      );
      throw AuthException(message);
    } catch (e, stackTrace) {
      dev.log(
        'Error signing in with Google',
        name: 'AuthService',
        error: e,
        stackTrace: stackTrace,
      );
      rethrow;
    }
  }

  Future<void> signInWithEmailAndPassword(String email, String password) async {
    try {
      await _auth.signInWithEmailAndPassword(email: email, password: password);
      await syncUser(); // trigger syncUser to backend
    } on auth.FirebaseAuthException catch (e, stackTrace) {
      final message = _handleFirebaseAuthError(e);
      dev.log(
        'Error signing in: $message',
        name: 'AuthService',
        error: e,
        stackTrace: stackTrace,
      );
      throw AuthException(message);
    } catch (e, stackTrace) {
      dev.log(
        'Error signing in with email and password',
        name: 'AuthService',
        error: e,
        stackTrace: stackTrace,
      );
      rethrow;
    }
  }

  Future<void> signUpWithEmailAndPassword(String email, String password) async {
    try {
      await _auth.createUserWithEmailAndPassword(
        email: email,
        password: password,
      );
      await syncUser(); // trigger syncUser to backend
    } on auth.FirebaseAuthException catch (e, stackTrace) {
      final message = _handleFirebaseAuthError(e);
      dev.log(
        'Error signing up: $message',
        name: 'AuthService',
        error: e,
        stackTrace: stackTrace,
      );
      throw AuthException(message);
    } catch (e, stackTrace) {
      dev.log(
        'Error signing up with email and password',
        name: 'AuthService',
        error: e,
        stackTrace: stackTrace,
      );
      rethrow;
    }
  }

  Future<void> signOut() async {
    await _auth.signOut();
  }

  String _handleFirebaseAuthError(auth.FirebaseAuthException e) {
    switch (e.code) {
      case 'weak-password':
        return 'The password provided is too weak.';
      case 'email-already-in-use':
        return 'An account already exists with this email.';
      case 'user-not-found':
        return 'No user found with this email.';
      case 'wrong-password':
        return 'Wrong password provided.';
      case 'invalid-email':
        return 'The email address is invalid.';
      case 'user-disabled':
        return 'This user account has been disabled.';
      case 'too-many-requests':
        return 'Too many requests. Please try again later.';
      case 'operation-not-allowed':
        return 'This operation is not allowed.';
      default:
        return 'Authentication error: ${e.message}';
    }
  }

  
}
