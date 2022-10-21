import 'dart:io';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:client/manage/oss.dart';
import 'package:flutter/material.dart';

/*
* 用户头像组件
* @author wuxubaiyang
* @Time 2022/10/21 10:28
*/
class Avatar extends StatelessWidget {
  // 头像地址
  final String uri;

  // 点击事件
  final VoidCallback? onTap;

  // 头像尺寸
  final AvatarSize avatarSize;

  // 头像来源
  final AvatarSource avatarSource;

  const Avatar({
    super.key,
    required this.uri,
    required this.avatarSource,
    this.onTap,
    this.avatarSize = AvatarSize.normal,
  });

  // oss资源头像
  const Avatar.oss({
    super.key,
    required this.uri,
    this.onTap,
    this.avatarSize = AvatarSize.normal,
  }) : avatarSource = AvatarSource.oss;

  // 网络头像
  const Avatar.net({
    super.key,
    required this.uri,
    this.onTap,
    this.avatarSize = AvatarSize.normal,
  }) : avatarSource = AvatarSource.net;

  // 本地头像
  const Avatar.file({
    super.key,
    required this.uri,
    this.onTap,
    this.avatarSize = AvatarSize.normal,
  }) : avatarSource = AvatarSource.file;

  // assets头像
  const Avatar.assets({
    super.key,
    required this.uri,
    this.onTap,
    this.avatarSize = AvatarSize.normal,
  }) : avatarSource = AvatarSource.assets;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: _buildAvatar(uri),
    );
  }

  // 构建头像
  Widget _buildAvatar(String uri) {
    if (uri.isEmpty) return _buildDefaultAvatar();
    switch (avatarSource) {
      case AvatarSource.net:
        return _buildNetAvatar(uri);
      case AvatarSource.oss:
        return _buildOSSAvatar(uri);
      case AvatarSource.file:
        return _buildFileAvatar(uri);
      case AvatarSource.assets:
        return _buildAssetsAvatar(uri);
    }
  }

  // 构建assets头像
  Widget _buildAssetsAvatar(String assetsName) {
    return CircleAvatar(
      radius: avatarSize.size,
      foregroundImage: AssetImage(assetsName),
      child: _buildDefaultAvatar(),
    );
  }

  // 构建本地头像
  Widget _buildFileAvatar(String path) {
    return CircleAvatar(
      radius: avatarSize.size,
      foregroundImage: FileImage(File(path)),
      child: _buildDefaultAvatar(),
    );
  }

  // 构建oss保管头像
  Widget _buildOSSAvatar(String objectKey) {
    return FutureBuilder<String>(
      future: ossManage.getObjectUrl(uri),
      builder: (_, snap) {
        return _buildNetAvatar(snap.data ?? "");
      },
    );
  }

  // 构建网络头像
  Widget _buildNetAvatar(String url) {
    return CachedNetworkImage(
      imageUrl: url,
      placeholder: (_, __) => _buildDefaultAvatar(),
      errorWidget: (_, __, err) => _buildDefaultAvatar(),
      imageBuilder: (_, image) {
        return CircleAvatar(
          radius: avatarSize.size,
          backgroundImage: image,
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
        AvatarSize.small: 14.0,
        AvatarSize.normal: 24.0,
        AvatarSize.large: 55.0,
      }[this]!;
}

// 头像来源枚举
enum AvatarSource { net, oss, file, assets }
