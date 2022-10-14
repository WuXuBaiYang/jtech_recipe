import 'package:client/common/manage.dart';
import 'package:client/manage/cache.dart';
import 'package:client/manage/event.dart';
import 'package:client/model/event.dart';
import 'package:flutter/material.dart';

/*
* 样式管理
* @author wuxubaiyang
* @Time 2022/10/14 10:09
*/
class ThemeManage extends BaseManage {
  // 默认样式缓存字段
  final String defaultThemeCacheKey = "default_theme_cache";

  // 是否使用Material3样式缓存字段
  final String defaultThemeUseMaterial3 = "default_theme_use_material3";

  static final ThemeManage _instance = ThemeManage._internal();

  factory ThemeManage() => _instance;

  ThemeManage._internal();

  @override
  Future<void> init() async {
    // 获取当前样式并发送消息切换样式
    eventManage.send(ThemeEvent(themeData: currentTheme));
  }

  // 当前样式
  ThemeData? get currentTheme {
    var useMaterial3 = cacheManage.getBool(defaultThemeUseMaterial3);
    var index = cacheManage.getInt(defaultThemeCacheKey);
    return ThemeType.values[index ?? 0].getTheme(useMaterial3);
  }

  // 切换默认样式
  Future<bool> switchTheme(ThemeType type, {bool? useMaterial3}) async {
    var result = await cacheManage.setInt(defaultThemeCacheKey, type.index);
    eventManage.send(ThemeEvent(
      themeData: type.getTheme(useMaterial3),
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
  blue,
  green,
}

// 样式枚举扩展
extension ThemeTypeExtension on ThemeType {
  // 样式中文名
  String? get nameCN => <ThemeType, String>{
        ThemeType.light: "日间模式",
        ThemeType.dark: "夜间模式",
        ThemeType.blue: "蓝色",
        ThemeType.green: "绿色",
      }[this];

  // 获取对应的样式配置
  ThemeData? getTheme(bool? useMaterial3) => <ThemeType, ThemeData>{
        ThemeType.light: ThemeData.light(useMaterial3: useMaterial3),
        ThemeType.dark: ThemeData.dark(useMaterial3: useMaterial3),
        ThemeType.blue: ThemeData(
          useMaterial3: useMaterial3,
          brightness: Brightness.light,
          colorSchemeSeed: Colors.lightBlueAccent,
        ),
        ThemeType.green: ThemeData(
          useMaterial3: useMaterial3,
          brightness: Brightness.light,
          colorSchemeSeed: Colors.lightGreenAccent,
        ),
      }[this];
}
