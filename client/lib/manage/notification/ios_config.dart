import 'package:client/common/model.dart';

/*
* IOS通知相关字段
* @author wuxubaiyang
* @Time 2021/8/31 2:06 下午
*/
class IOSNotificationConfig extends BaseModel {
  // 是否通知
  final bool? presentAlert;

  // 是否标记
  final bool? presentBadge;

  // 是否有声音
  final bool? presentSound;

  // 声音文件
  final String? sound;

  // 标记数字
  final int? badgeNumber;

  // 子标题
  final String? subtitle;

  // 线程标识
  final String? threadIdentifier;

  const IOSNotificationConfig({
    this.presentAlert,
    this.presentBadge,
    this.presentSound,
    this.sound,
    this.badgeNumber,
    this.subtitle,
    this.threadIdentifier,
  });
}
