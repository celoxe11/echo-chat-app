import 'package:echo_chat_app_frontend/core/routes/app_routes.dart';
import 'package:echo_chat_app_frontend/presentation/screens/login_screen.dart';
import 'package:echo_chat_app_frontend/presentation/screens/register_screen.dart';
import 'package:flutter/material.dart';

class AppRouter {
  Route? onGenerateRoute(RouteSettings settings) {
    switch (settings.name) {
      case AppRoutes.login:
        return _buildRoute(const LoginScreen());
      case AppRoutes.register:
        return _buildRoute(const RegisterScreen());

      // case AppRoutes.chat:
      //   return _buildRoute(
      //     BlocProvider(
      //       // Ini BLoC lokal (hanya untuk halaman chat)
      //       create: (context) => ChatBloc(
      //         chatRepository: context.read<ChatRepository>(),
      //       ),
      //       child: const ChatScreen(),
      //     ),
      //   );

      default:
        return _buildRoute(const Scaffold(body: Center(child: Text("404"))));
    }
  }

  MaterialPageRoute _buildRoute(Widget child) {
    return MaterialPageRoute(builder: (_) => child);
  }
}
