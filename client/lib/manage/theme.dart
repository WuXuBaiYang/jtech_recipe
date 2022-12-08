import 'package:client/common/manage.dart';
import 'package:client/manage/cache.dart';
import 'package:client/manage/event.dart';
import 'package:flutter/material.dart';

/*
* 样式管理
* @author wuxubaiyang
* @Time 2022/10/14 10:09
*/
class ThemeManage extends BaseManage {
  // 默认样式缓存字段
  final String _defaultThemeCacheKey = 'default_theme_cache';

  static final ThemeManage _instance = ThemeManage._internal();

  factory ThemeManage() => _instance;

  ThemeManage._internal();

  // 获取主色
  Color get primaryColor => currentTheme.primaryColor;

  // 当前样式
  ThemeData get currentTheme {
    final index = cacheManage.getInt(_defaultThemeCacheKey);
    return ThemeType.values[index ?? 0].theme;
  }

  // 切换默认样式
  Future<bool> switchTheme(ThemeType type) async {
    final result = await cacheManage.setInt(_defaultThemeCacheKey, type.index);
    eventManage.send(ThemeEvent(
      themeData: type.theme,
    ));
    return result;
  }
}

// 单例调用
final themeManage = ThemeManage();

// 支持的样式枚举
enum ThemeType {
  light,
  dark,
}

// 样式枚举扩展
extension ThemeTypeExtension on ThemeType {
  // 样式中文名
  String get nameCN => <ThemeType, String>{
        ThemeType.light: '日间模式',
        ThemeType.dark: '夜间模式',
      }[this]!;

  // 获取对应的样式配置
  ThemeData get theme => <ThemeType, ThemeData>{
        ThemeType.light: ThemeData(
          useMaterial3: true,
          brightness: Brightness.light,
          inputDecorationTheme: _inputDecorationTheme,
        ),
        ThemeType.dark: ThemeData(
          useMaterial3: true,
          brightness: Brightness.dark,
          inputDecorationTheme: _inputDecorationTheme,
        ),
      }[this]!;

  // 获取输入框样式
  InputDecorationTheme get _inputDecorationTheme => const InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.all(
            Radius.circular(6),
          ),
        ),
        floatingLabelBehavior: FloatingLabelBehavior.always,
        contentPadding: EdgeInsets.all(12),
        alignLabelWithHint: true,
      );
}

/*
* 全局样式控制事件
* @author wuxubaiyang
* @Time 2022/4/1 15:14
*/
class ThemeEvent extends EventModel {
  // 全局样式
  final ThemeData themeData;

  const ThemeEvent({required this.themeData});
}
