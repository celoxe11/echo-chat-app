import 'package:echo_chat_app_frontend/core/api/dio_client.dart';
import 'package:echo_chat_app_frontend/data/models/user_model.dart';

class UserService {
  final DioClient dioClient = DioClient.instance;

  Future<User> getMe() async {
    final response = await dioClient.get('/users/me');
    return User.fromJson(response.data);
  }

  Future<User> searchUserByUsername(String username) async {
    final response = await dioClient.get('/users/search?username=$username');
    return User.fromJson(response.data);
  }

  Future<User> updateProfile(String username, String displayName, String? photoUrl) async {
    final response = await dioClient.patch('/users/me', data: {
      'username': username,
      'displayName': displayName,
      'photoUrl': photoUrl,
    });
    return User.fromJson(response.data);
  }
}