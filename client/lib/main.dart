import 'dart:io';

import 'package:client/common/common.dart';
import 'package:client/common/logic.dart';
import 'package:client/manage/api_cancel.dart';
import 'package:client/manage/cache.dart';
import 'package:client/manage/event.dart';
import 'package:client/manage/notification/notification.dart';
import 'package:client/manage/oss.dart';
import 'package:client/manage/router.dart';
import 'package:client/manage/tag.dart';
import 'package:client/manage/theme.dart';
import 'package:client/model/event.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

import 'common/localization/chinese_cupertino_localizations.dart';
import 'manage/auth.dart';

// 调试状态
bool debugMode = true;

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  // 初始化各种manage
  await routerManage.init(); // 路由服务
  await noticeManage.init(); // 通知服务
  await apiCancelManage.init(); // 请求撤销服务
  await cacheManage.init(); // 缓存服务
  await eventManage.init(); // 事件服务
  await authManage.init(); // 初始化授权服务
  await ossManage.init(); // oss服务
  await themeManage.init(); // 初始化样式管理
  await tagManage.init(); // 初始化标签管理
  // 设置沉浸式状态栏
  if (Platform.isAndroid) {
    SystemChrome.setSystemUIOverlayStyle(SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness:
          themeManage.currentTheme.brightness == Brightness.light
              ? Brightness.dark
              : Brightness.light,
    ));
  }
  runApp(const MyApp());
}

/*
* App根节点
* @author wuxubaiyang
* @Time 2022/9/8 14:47
*/
class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return StreamBuilder<ThemeEvent>(
      stream: eventManage.on<ThemeEvent>(),
      builder: (c, snap) => MaterialApp(
        title: 'JTech Recipe',
        navigatorKey: routerManage.navigateKey,
        debugShowCheckedModeBanner: debugMode,
        theme: snap.data?.themeData ?? themeManage.currentTheme,
        onGenerateRoute: routerManage.onGenerateRoute(
          routesMap: RoutePath.routes,
        ),
        localizationsDelegates: const [
          GlobalWidgetsLocalizations.delegate,
          GlobalMaterialLocalizations.delegate,
          ChineseCupertinoLocalizations.delegate,
        ],
        supportedLocales: const [
          Locale('en', 'US'),
          Locale('zh', 'CN'),
        ],
        home: const SplashPage(),
      ),
    );
  }
}

/*
* 欢迎页面
* @author wuxubaiyang
* @Time 2022/9/8 14:47
*/
class SplashPage extends StatefulWidget {
  const SplashPage({super.key});

  @override
  State<StatefulWidget> createState() => _SplashPageState();
}

/*
* 欢迎页-状态
* @author wuxubaiyang
* @Time 2022/9/8 14:48
*/
class _SplashPageState extends State<SplashPage> {
  // 逻辑管理
  final _logic = _SplashPageLogic();

  @override
  void initState() {
    super.initState();
    // 延迟跳转到下一页
    _logic.goNextPageDelay();
  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      body: Center(
        child: FlutterLogo(
          size: 55,
        ),
      ),
    );
  }

  @override
  void dispose() {
    _logic.dispose();
    super.dispose();
  }
}

/*
* 欢迎页-逻辑
* @author wuxubaiyang
* @Time 2022/10/27 11:15
*/
class _SplashPageLogic extends BaseLogic {
  // 延迟跳转到下一页
  Future<void> goNextPageDelay() async {
    await Future.delayed(const Duration(milliseconds: 600));
    if (authManage.authorized) {
      // 已授权则跳转到首页
      routerManage.pushReplacementNamed(RoutePath.home);
    } else {
      // 未授权跳转到授权页
      routerManage.pushReplacementNamed(RoutePath.auth);
    }
  }
}
