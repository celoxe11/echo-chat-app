class User {
  final int id;
  final String firebaseUid;
  final String email;
  final String fullName;
  final String username;
  final String avatarUrl;
  final String status;
  final DateTime lastSeen;
  final DateTime createdAt;
  final DateTime updatedAt;

  User({
    required this.id,
    required this.firebaseUid,
    required this.email,
    required this.fullName,
    required this.username,
    required this.avatarUrl,
    required this.status,
    required this.lastSeen,
    required this.createdAt,
    required this.updatedAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      firebaseUid: json['firebase_uid'],
      email: json['email'],
      fullName: json['full_name'],
      username: json['username'],
      avatarUrl: json['avatar_url'] ?? '',
      status: json['status'] ?? 'offline',
      lastSeen: DateTime.parse(json['last_seen']),
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }
}
