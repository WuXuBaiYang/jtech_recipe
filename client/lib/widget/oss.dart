import 'package:client/manage/oss.dart';
import 'package:client/widget/avatar.dart';
import 'package:client/widget/image.dart';
import 'package:flutter/widgets.dart';

/*
* oss图片
* @author wuxubaiyang
* @Time 2022/10/24 10:49
*/
class OSSImageView extends StatelessWidget {
  // oss的对象key
  final String object;

  // 宽，高，正方体尺寸
  final double? width, height, squareSize;

  // 图片填充方式
  final BoxFit? fit;

  // 图片形状
  final BoxShape? shape;

  // 图片边框
  final BoxBorder? border;

  // 边框圆角
  final BorderRadius? borderRadius;

  // 图片点击事件
  final VoidCallback? onTap;

  // 图片异常状态
  final ImageViewStateLoad? failState;

  // 图片加载中状态
  final ImageViewStateLoad? loadingState;

  // 图片加载完成状态
  final ImageViewStateLoad? completedState;

  const OSSImageView({
    super.key,
    required this.object,
    this.width,
    this.height,
    this.squareSize,
    this.fit,
    this.shape,
    this.border,
    this.borderRadius,
    this.onTap,
    this.failState,
    this.loadingState,
    this.completedState,
  });

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<String>(
      future: ossManage.getObjectUrl(object),
      builder: (_, snap) {
        if (snap.data?.isNotEmpty ?? false) {
          return ImageView.net(
            url: snap.data ?? '',
            width: width,
            height: height,
            squareSize: squareSize,
            fit: fit,
            shape: shape,
            border: border,
            borderRadius: borderRadius,
            onTap: onTap,
            failState: failState,
            loadingState: loadingState,
            completedState: completedState,
          );
        }
        return const SizedBox();
      },
    );
  }
}

/*
* oss头像
* @author wuxubaiyang
* @Time 2022/10/24 11:13
*/
class OSSAvatar extends StatelessWidget {
  // 头像尺寸
  final AvatarSize avatarSize;

  // oss的对象key
  final String object;

  // 图片点击事件
  final VoidCallback? onTap;

  const OSSAvatar({
    super.key,
    required this.object,
    this.onTap,
    this.avatarSize = AvatarSize.normal,
  });

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<String>(
      future: ossManage.getObjectUrl(object),
      builder: (_, snap) {
        return Avatar.net(
          url: snap.data ?? '',
          avatarSize: avatarSize,
          onTap: onTap,
        );
      },
    );
  }
}
