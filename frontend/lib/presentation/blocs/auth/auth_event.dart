import 'package:echo_chat_app_frontend/data/models/user_model.dart';
import 'package:equatable/equatable.dart';

abstract class AuthEvent extends Equatable {
  @override
  List<Object?> get props => [];
}

class AppStarted extends AuthEvent {
  AppStarted();
}

class AuthStatusChanged extends AuthEvent {
  final User? user;
  AuthStatusChanged(this.user);

  @override
  List<Object?> get props => [user];
}

class GoogleSignInRequested extends AuthEvent {
  GoogleSignInRequested();
}

class EmailSignInRequested extends AuthEvent {
  final String email;
  final String password;
  EmailSignInRequested({required this.email, required this.password});

  @override
  List<Object?> get props => [email, password];
}

class EmailSignUpRequested extends AuthEvent {
  final String email;
  final String password;
  EmailSignUpRequested({required this.email, required this.password});

  @override
  List<Object?> get props => [email, password];
}

class SignOutRequested extends AuthEvent {
  SignOutRequested();
}

class ForgotPasswordRequested extends AuthEvent {
  final String email;
  ForgotPasswordRequested({required this.email});

  @override
  List<Object?> get props => [email];
}