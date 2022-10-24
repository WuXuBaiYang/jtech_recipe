import 'dart:io';
import 'dart:typed_data';

import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';

import 'image.dart';

/*
* 用户头像组件
* @author wuxubaiyang
* @Time 2022/10/21 10:28
*/
class Avatar extends StatelessWidget {
  // 头像尺寸
  final AvatarSize avatarSize;

  // 图片源
  final ImageViewSource source;

  // 图片点击事件
  final VoidCallback? onTap;

  const Avatar({
    super.key,
    required this.source,
    this.onTap,
    this.avatarSize = AvatarSize.normal,
  });

  // 本地头像
  Avatar.file({
    super.key,
    required File file,
    this.onTap,
    bool? cacheRawData,
    String? imageCacheName,
    this.avatarSize = AvatarSize.normal,
  }) : source = ImageViewSource.file(
          file: file,
          cacheRawData: cacheRawData,
          imageCacheName: imageCacheName,
        );

  // assets头像
  Avatar.assets({
    super.key,
    required String assetName,
    this.onTap,
    AssetBundle? bundle,
    String? package,
    bool? cacheRawData,
    String? imageCacheName,
    this.avatarSize = AvatarSize.normal,
  }) : source = ImageViewSource.assets(
          assetName: assetName,
          bundle: bundle,
          package: package,
          cacheRawData: cacheRawData,
          imageCacheName: imageCacheName,
        );

  // 内存头像
  Avatar.memory({
    super.key,
    required Uint8List bytes,
    this.onTap,
    bool? cacheRawData,
    String? imageCacheName,
    this.avatarSize = AvatarSize.normal,
  }) : source = ImageViewSource.memory(
          bytes: bytes,
          cacheRawData: cacheRawData,
          imageCacheName: imageCacheName,
        );

  // 网络头像
  Avatar.net({
    super.key,
    required String url,
    this.onTap,
    Map<String, String>? headers,
    bool? cache,
    int? retries,
    Duration? timeLimit,
    Duration? timeRetry,
    CancellationToken? cancelToken,
    String? cacheKey,
    bool? cacheRawData,
    String? imageCacheName,
    this.avatarSize = AvatarSize.normal,
  }) : source = ImageViewSource.net(
          url: url,
          headers: headers,
          cache: cache,
          retries: retries,
          timeLimit: timeLimit,
          timeRetry: timeRetry,
          cancelToken: cancelToken,
          cacheKey: cacheKey,
          cacheRawData: cacheRawData,
          imageCacheName: imageCacheName,
        );

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: CircleAvatar(
        radius: avatarSize.size,
        foregroundImage: source.provider,
        child: _buildDefaultAvatar(),
      ),
    );
  }

  // 构建默认头像
  Widget _buildDefaultAvatar() {
    return Icon(
      Icons.account_circle_rounded,
      size: avatarSize.defSize,
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
        AvatarSize.large: 45.0,
      }[this]!;

  // 获取默认头像尺寸
  double get defSize => {
        AvatarSize.small: 25.0,
        AvatarSize.normal: 45.0,
        AvatarSize.large: 55.0,
      }[this]!;
}
