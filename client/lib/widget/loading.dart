import 'package:client/manage/theme.dart';
import 'package:flutter/material.dart';

/*
* 加载状态组件
* @author wuxubaiyang
* @Time 2022/10/20 10:32
*/
class Loading extends StatelessWidget {
  // 默认展示元素
  final Widget child;

  // 是否加载中
  final bool loading;

  // 加载组件尺寸
  final LoadingSize loadingSize;

  // 加载组件样式
  final LoadingStyle loadingStyle;

  const Loading({
    super.key,
    required this.loading,
    required this.child,
    this.loadingSize = LoadingSize.normal,
    this.loadingStyle = LoadingStyle.light,
  });

  // 深色构造
  const Loading.dark({
    super.key,
    required this.loading,
    required this.child,
    this.loadingSize = LoadingSize.normal,
  }) : loadingStyle = LoadingStyle.dark;

  @override
  Widget build(BuildContext context) {
    if (loading) {
      return SizedBox.square(
        dimension: loadingSize.size,
        child: CircularProgressIndicator(
          color: loadingStyle.color,
        ),
      );
    }
    return child;
  }
}

// 加载组件尺寸
enum LoadingSize { small, normal, large }

// 加载组件尺寸扩展
extension LoadingSizeExtension on LoadingSize {
  // 获取真实尺寸
  double get size => {
        LoadingSize.small: 14.0,
        LoadingSize.normal: 18.0,
        LoadingSize.large: 22.0,
      }[this]!;
}

// 加载组件样式
enum LoadingStyle { light, dark }

// 加载组件样式扩展
extension LoadingStyleExtension on LoadingStyle {
  // 获取样式颜色
  Color get color => {
        LoadingStyle.light: themeManage.currentTheme.primaryColor,
        LoadingStyle.dark: Colors.white,
      }[this]!;
}
