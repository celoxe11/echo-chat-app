import 'package:echo_chat_app_frontend/data/models/user_model.dart';
import 'package:echo_chat_app_frontend/data/services/user_service.dart';

class UserRepository {
  final UserService userService;

  UserRepository({required this.userService});

  Future<User> getMe() async {
    return await userService.getMe();
  }

  Future<User> searchUserByUsername(String username) async {
    return await userService.searchUserByUsername(username);
  }

  Future<User> updateProfile(String username, String displayName, String? photoUrl) async {
    return await userService.updateProfile(username, displayName, photoUrl);
  }
}