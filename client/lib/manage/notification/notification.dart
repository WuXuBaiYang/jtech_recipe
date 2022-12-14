import 'dart:io';
import 'package:client/common/manage.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'android_config.dart';
import 'ios_config.dart';

// 当收到消息通知时的回调
typedef OnNotificationReceive = Future Function(
    int id, String? title, String? body, String? payload);

// 当通知消息被点击触发时
typedef OnNotificationSelect = Future Function(String? payload);

/*
* 本地通知管理
* @author JTech JH
* @Time 2022/3/29 10:46
*/
class NotificationManage extends BaseManage {
  static final NotificationManage _instance = NotificationManage._internal();

  factory NotificationManage() => _instance;

  NotificationManage._internal();

  // 默认图标名称
  final String _defaultIconName = "ic_launcher";

  // 接受通知消息回调集合
  final List<OnNotificationReceive> _notificationReceiveListeners = [];

  // 通知消息点击触发回调集合
  final List<OnNotificationSelect> _notificationSelectListeners = [];

  // 通知推送管理
  late FlutterLocalNotificationsPlugin localNotification;

  // 通知栏初始化状态记录
  bool? _initialized;

  @override
  Future<void> init() async {
    localNotification = FlutterLocalNotificationsPlugin();
  }

  // 获取初始化状态
  bool get initialized => _initialized ?? false;

  // 初始化通知栏消息
  Future<bool?> initNotification(String icon) => localNotification
      .initialize(
        InitializationSettings(
          android: AndroidInitializationSettings(icon),
          iOS: IOSInitializationSettings(
            onDidReceiveLocalNotification: _onReceiveNotification,
          ),
        ),
        onSelectNotification: _onNotificationSelect,
      )
      .then((value) => _initialized = value);

  // 显示进度通知
  Future<void> showProgress({
    required int maxProgress,
    required int progress,
    required bool indeterminate,
    // 基础参数
    required int id,
    String? title,
    String? body,
    String? payload,
  }) {
    if (null == body && !indeterminate) {
      double ratio = (progress / maxProgress.toDouble()) * 100;
      body = "${ratio.toStringAsFixed(1)}%";
    }
    return show(
      id: id,
      title: title,
      body: body,
      payload: payload,
      androidConfig: AndroidNotificationConfig(
        showProgress: true,
        maxProgress: maxProgress,
        progress: progress,
        indeterminate: indeterminate,
        playSound: false,
        enableLights: false,
        enableVibration: false,
        ongoing: true,
        onlyAlertOnce: true,
      ),
      iosConfig: const IOSNotificationConfig(
        presentSound: false,
        presentBadge: false,
      ),
    );
  }

  // 显示通知栏消息
  Future<void> show({
    required int id,
    String? title,
    String? body,
    String? payload,
    AndroidNotificationConfig? androidConfig,
    IOSNotificationConfig? iosConfig,
  }) async {
    if (!initialized) await initNotification(_defaultIconName);
    assert(
        initialized,
        "请在 android-app-src-main-res-drawable 目录下添加 app_icon 图片文件；"
        "或者调用 jNotificationManage.initNotification() 自行指定默认图标");
    // 申请ios权限
    if (Platform.isIOS) {
      var result = await localNotification
          .resolvePlatformSpecificImplementation<
              IOSFlutterLocalNotificationsPlugin>()
          ?.requestPermissions(alert: true, badge: true, sound: true);
      if (null != result && !result) return;
    }
    androidConfig ??= const AndroidNotificationConfig();
    iosConfig ??= const IOSNotificationConfig();
    return localNotification.show(
      id,
      title,
      body,
      NotificationDetails(
        android: AndroidNotificationDetails(
          androidConfig.channelId ?? "$id",
          androidConfig.channelName ?? "$id",
          channelDescription: androidConfig.channelDescription ?? "$id",
          channelShowBadge: androidConfig.channelShowBadge,
          importance: Importance.max,
          priority: Priority.high,
          showWhen: null != androidConfig.when,
          when: androidConfig.when?.inMilliseconds ?? 0,
          icon: androidConfig.icon,
          playSound: androidConfig.playSound,
          enableVibration: androidConfig.enableVibration,
          groupKey: androidConfig.groupKey,
          setAsGroupSummary: androidConfig.setAsGroupSummary,
          autoCancel: androidConfig.autoCancel,
          ongoing: androidConfig.ongoing,
          onlyAlertOnce: androidConfig.onlyAlertOnce,
          enableLights: androidConfig.enableLights,
          timeoutAfter: androidConfig.timeoutAfter?.inMilliseconds,
          showProgress: androidConfig.showProgress,
          maxProgress: androidConfig.maxProgress,
          progress: androidConfig.progress,
          indeterminate: androidConfig.indeterminate,
        ),
        iOS: IOSNotificationDetails(
          presentAlert: iosConfig.presentAlert,
          presentBadge: iosConfig.presentBadge,
          presentSound: iosConfig.presentSound,
          sound: iosConfig.sound,
          badgeNumber: iosConfig.badgeNumber,
          subtitle: iosConfig.subtitle,
          threadIdentifier: iosConfig.threadIdentifier,
        ),
      ),
      payload: payload,
    );
  }

  // 取消通知
  Future<void> cancel(int id, {String? tag}) =>
      localNotification.cancel(id, tag: tag);

  // 取消所有通知
  Future<void> cancelAll() => localNotification.cancelAll();

  // 添加接受消息监听
  void addReceiveListener(OnNotificationReceive listener) =>
      _notificationReceiveListeners.add(listener);

  // 添加消息选择监听
  void addSelectListener(OnNotificationSelect listener) =>
      _notificationSelectListeners.add(listener);

  // 当接收到通知消息回调
  Future _onReceiveNotification(
      int id, String? title, String? body, String? payload) async {
    for (var item in _notificationReceiveListeners) {
      await item(id, title, body, payload);
    }
  }

  // 消息通知点击事件回调
  Future _onNotificationSelect(String? payload) async {
    for (var item in _notificationSelectListeners) {
      await item(payload);
    }
  }
}

// 单例调用
final noticeManage = NotificationManage();
