import 'package:cached_network_image/cached_network_image.dart';
import 'package:client/manage/oss.dart';
import 'package:client/model/user.dart';
import 'package:flutter/material.dart';

/*
* 用户头像组件
* @author wuxubaiyang
* @Time 2022/10/21 10:28
*/
class Avatar extends StatelessWidget {
  // 用户对象
  final UserModel user;

  // 头像尺寸
  final AvatarSize avatarSize;

  const Avatar({
    super.key,
    required this.user,
    this.avatarSize = AvatarSize.normal,
  });

  // 小号头像
  const Avatar.small({
    super.key,
    required this.user,
  }) : avatarSize = AvatarSize.small;

  // 大号头像
  const Avatar.large({
    super.key,
    required this.user,
  }) : avatarSize = AvatarSize.large;

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<String>(
      future: ossManage.getObjectUrl(user.avatar),
      builder: (_, snap) {
        return CachedNetworkImage(
          imageUrl: snap.data ?? "",
          placeholder: (_, __) => _buildDefaultAvatar(),
          errorWidget: (_, __, err) => _buildDefaultAvatar(),
          imageBuilder: (_, image) {
            return CircleAvatar(
              radius: avatarSize.size,
              backgroundImage: image,
            );
          },
        );
      },
    );
  }

  // 构建默认头像
  Widget _buildDefaultAvatar() {
    return CircleAvatar(
      radius: avatarSize.size,
      child: Icon(
        Icons.account_circle_rounded,
        size: avatarSize.size * 2 - 8,
      ),
    );
  }
}

// 头像组件尺寸枚举
enum AvatarSize { small, normal, large }

// 头像组件尺寸扩展
extension AvatarSizeExtension on AvatarSize {
  // 获取尺寸
  double get size => {
        AvatarSize.small: 18.0,
        AvatarSize.normal: 24.0,
        AvatarSize.large: 32.0,
      }[this]!;
}
