import 'package:echo_chat_app_frontend/core/routes/app_routes.dart';
import 'package:echo_chat_app_frontend/data/repositories/auth_repository.dart';
import 'package:echo_chat_app_frontend/data/repositories/user_repository.dart';
import 'package:echo_chat_app_frontend/data/services/auth_service.dart';
import 'package:echo_chat_app_frontend/data/services/user_service.dart';
import 'package:echo_chat_app_frontend/presentation/blocs/auth/auth_bloc.dart';
import 'package:echo_chat_app_frontend/presentation/blocs/auth/auth_event.dart';
import 'package:echo_chat_app_frontend/presentation/blocs/auth/auth_state.dart';
import 'package:flutter/material.dart';
import 'package:echo_chat_app_frontend/core/config/firebase_config.dart';
import 'package:echo_chat_app_frontend/core/routes/app_router.dart';
import 'package:echo_chat_app_frontend/core/theme/app_theme.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await FirebaseConfig.initialize();

  // initiate all services and repository here
  final authService = AuthService();
  final authRepository = AuthRepository(authService: authService);

  final userService = UserService();
  final userRepository = UserRepository(userService: userService);

  runApp(
    MultiRepositoryProvider(
      providers: [
        // provide all repository here
        RepositoryProvider<AuthRepository>.value(value: authRepository),
        RepositoryProvider<UserRepository>.value(value: userRepository),
      ],
      child: MultiBlocProvider(
        providers: [
          // global bloc provider
          BlocProvider<AuthBloc>(
            create: (context) =>
                AuthBloc(authRepository: context.read<AuthRepository>())
                  ..add(AppStarted()), // Event awal untuk cek login status
          ),
        ],
        child: const MyApp(),
      ),
    ),
  );
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Echo',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      onGenerateRoute: AppRouter().onGenerateRoute,
      builder: (context, child) {
        return BlocListener<AuthBloc, AuthState>(
          listener: (context, state) {
            if (state is AuthAuthenticated) {
              // Jika login sukses, pindah ke Home
              Navigator.of(
                context,
              ).pushNamedAndRemoveUntil(AppRoutes.chats, (route) => false);
            } else if (state is AuthUnauthenticated) {
              // Jika logout/belum login, pindah ke Login
              Navigator.of(
                context,
              ).pushNamedAndRemoveUntil(AppRoutes.login, (route) => false);
            }
          },
          child: child,
        );
      },
      initialRoute: AppRoutes.login,
    );
  }
}
