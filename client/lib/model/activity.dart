import 'package:client/common/model.dart';
import 'package:client/model/model.dart';
import 'package:client/tool/date.dart';

/*
* 活动信息
* @author wuxubaiyang
* @Time 2022/10/14 14:56
*/
class ActivityModel extends BaseModel with BasePart {
  // 持续周期
  final num cycleTime;

  // 是否为长期活动
  final bool always;

  // 活动标题
  final String title;

  // 介绍网址
  final String url;

  // 支持的活动类型
  final List<String> typeCodes;

  ActivityModel.from(obj)
      : cycleTime = obj?['cycleTime'] ?? 0,
        always = obj?['always'] ?? false,
        title = obj?['title'] ?? '',
        url = obj?['url'] ?? '',
        typeCodes =
            (obj?['typeCodes'] ?? []).map<String>((e) => '$e').toList() {
    initBasePart(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        'cycleTime': cycleTime,
        'always': always,
        'title': title,
        'url': url,
        'typeCodes': typeCodes,
      };
}

/*
* 活动记录
* @author wuxubaiyang
* @Time 2022/10/14 14:59
*/
class ActivityRecordModel extends BaseModel with BasePart {
  // 开始时间
  final DateTime beginTime;

  // 结束时间
  final DateTime endTime;

  // 活动id
  final String activityId;

  // 活动信息
  final ActivityModel activity;

  ActivityRecordModel.from(obj)
      : beginTime = DateTime.tryParse(obj?['beginTime'] ?? '') ?? DateTime(0),
        endTime = DateTime.tryParse(obj?['endTime'] ?? '') ?? DateTime(0),
        activityId = obj?['activityId'] ?? '',
        activity = ActivityModel.from(obj?['activity'] ?? {}) {
    initBasePart(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        'beginTime': beginTime.toIso8601StringWithUTC(),
        'endTime': endTime.toIso8601StringWithUTC(),
        'activityId': activityId,
        'activity': activity.to(),
      };
}
