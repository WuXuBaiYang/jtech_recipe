import 'dart:io';
import 'dart:typed_data';

import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';

// 图片状态加载
typedef ImageViewStateLoad = Widget? Function();

/*
* 图片
* @author wuxubaiyang
* @Time 2022/10/24 9:54
*/
class ImageView extends StatefulWidget {
  // 图片源
  final ImageViewSource source;

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

  const ImageView({
    super.key,
    required this.source,
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

  // 本地图片
  ImageView.file({
    super.key,
    required File file,
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
    bool? cacheRawData,
    String? imageCacheName,
  }) : source = ImageViewSource.file(
          file: file,
          cacheRawData: cacheRawData,
          imageCacheName: imageCacheName,
        );

  // assets图片
  ImageView.assets({
    super.key,
    required String assetName,
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
    AssetBundle? bundle,
    String? package,
    bool? cacheRawData,
    String? imageCacheName,
  }) : source = ImageViewSource.assets(
          assetName: assetName,
          bundle: bundle,
          package: package,
          cacheRawData: cacheRawData,
          imageCacheName: imageCacheName,
        );

  // 内存图片
  ImageView.memory({
    super.key,
    required Uint8List bytes,
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
    bool? cacheRawData,
    String? imageCacheName,
  }) : source = ImageViewSource.memory(
          bytes: bytes,
          cacheRawData: cacheRawData,
          imageCacheName: imageCacheName,
        );

  // 网络图片
  ImageView.net({
    super.key,
    required String url,
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
    Map<String, String>? headers,
    bool? cache,
    int? retries,
    Duration? timeLimit,
    Duration? timeRetry,
    CancellationToken? cancelToken,
    String? cacheKey,
    bool? cacheRawData,
    String? imageCacheName,
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
  State<StatefulWidget> createState() => _ImageViewState();
}

/*
* 图片-状态
* @author wuxubaiyang
* @Time 2022/10/24 9:54
*/
class _ImageViewState extends State<ImageView> {
  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: widget.onTap,
      child: ExtendedImage(
        image: widget.source.provider,
        fit: widget.fit,
        shape: widget.shape,
        border: widget.border,
        enableLoadState: true,
        borderRadius: widget.borderRadius,
        width: widget.squareSize ?? widget.width,
        height: widget.squareSize ?? widget.height,
        loadStateChanged: _loadStateChanged,
      ),
    );
  }

  // 图片加载状态
  Widget? _loadStateChanged(ExtendedImageState state) {
    switch (state.extendedImageLoadState) {
      case LoadState.loading:
        return _buildLoadingState(state);
      case LoadState.failed:
        return _buildLoadFailState(state);
      case LoadState.completed:
        return widget.completedState?.call() ?? state.completedWidget;
    }
  }

  // 构建加载中状态
  Widget _buildLoadingState(ExtendedImageState state) {
    return widget.loadingState?.call() ??
        const Center(
          child: CircularProgressIndicator(),
        );
  }

  // 构建加载失败状态
  Widget _buildLoadFailState(ExtendedImageState state) {
    return widget.failState?.call() ??
        GestureDetector(
          child: Center(
            child: Icon(
              Icons.refresh_rounded,
              color: Colors.grey[400],
            ),
          ),
          onTap: () => state.reLoadImage(),
        );
  }
}

/*
* 图片数据源
* @author wuxubaiyang
* @Time 2022/10/24 9:56
*/
class ImageViewSource {
  // 图片代理
  final ImageProvider provider;

  ImageViewSource({
    required this.provider,
  });

  // 本地图片
  ImageViewSource.file({
    required File file,
    bool? cacheRawData,
    String? imageCacheName,
  }) : provider = ExtendedFileImageProvider(
          file,
          cacheRawData: cacheRawData ?? false,
          imageCacheName: imageCacheName,
        );

  // assets图片
  ImageViewSource.assets({
    required String assetName,
    AssetBundle? bundle,
    String? package,
    bool? cacheRawData,
    String? imageCacheName,
  }) : provider = ExtendedAssetImageProvider(
          assetName,
          bundle: bundle,
          package: package,
          cacheRawData: cacheRawData ?? false,
          imageCacheName: imageCacheName,
        );

  // 内存图片
  ImageViewSource.memory({
    required Uint8List bytes,
    bool? cacheRawData,
    String? imageCacheName,
  }) : provider = ExtendedMemoryImageProvider(
          bytes,
          cacheRawData: cacheRawData ?? false,
          imageCacheName: imageCacheName,
        );

  // 网络图片
  ImageViewSource.net({
    required String url,
    Map<String, String>? headers,
    bool? cache,
    int? retries,
    Duration? timeLimit,
    Duration? timeRetry,
    CancellationToken? cancelToken,
    String? cacheKey,
    bool? cacheRawData,
    String? imageCacheName,
  }) : provider = ExtendedNetworkImageProvider(
          url,
          headers: headers,
          cache: cache ?? true,
          retries: retries ?? 1,
          timeLimit: timeLimit,
          timeRetry: timeRetry ?? const Duration(milliseconds: 1000),
          cancelToken: cancelToken,
          cacheKey: cacheKey,
          cacheRawData: cacheRawData ?? false,
          imageCacheName: imageCacheName,
        );
}