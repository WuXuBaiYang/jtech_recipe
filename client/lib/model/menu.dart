import 'package:client/common/model.dart';
import 'package:client/model/activity.dart';
import 'package:client/model/model.dart';

/*
* 菜单信息
* @author wuxubaiyang
* @Time 2022/10/14 15:16
*/
class MenuModel extends BaseModel with BasePart, CreatorPart {
  // 内容集合
  final List<MenuContentItem> contents;

  // 来源id
  final String? originId;

  // 来源菜单
  final MenuModel? originMenu;

  // 活动记录id
  final String? activityRecordId;

  // 活动记录信息
  final ActivityRecordModel? activityRecord;

  // 是否点赞
  final bool liked;

  // 点赞数量
  final num likeCount;

  // 是否收藏
  final bool collected;

  // 收藏数量
  final num collectCount;

  MenuModel.from(obj)
      : contents = (obj?["contents"] ?? [])
            .map<MenuContentItem>((e) => MenuContentItem.from(e))
            .toList(),
        originId = obj?["originId"],
        originMenu = obj?["originMenu"] != null
            ? MenuModel.from(obj?["originMenu"] ?? {})
            : null,
        activityRecordId = obj?["activityRecordId"],
        activityRecord = obj?["activityRecord"] != null
            ? ActivityRecordModel.from(obj?["activityRecord"] ?? {})
            : null,
        liked = obj?["liked"] ?? false,
        likeCount = obj?["likeCount"] ?? 0,
        collected = obj?["collected"] ?? false,
        collectCount = obj?["collectCount"] ?? 0 {
    initBasePart(obj);
    initCreatorPart(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        ...creatorPart,
        "contents": contents.map((e) => e.to()).toList(),
        "originId": originId,
        "originMenu": originMenu?.to(),
        "activityRecordId": activityRecordId,
        "activityRecord": activityRecord?.to(),
        "liked": liked,
        "likeCount": likeCount,
        "collected": collected,
        "collectCount": collectCount,
      };

  // 获取编辑结构
  Map<String, dynamic> toModifyInfo() => {
        "contents": contents.map((e) => e.to()).toList(),
        "activityRecordId": activityRecordId,
      };
}

/*
* 菜单内容结构
* @author wuxubaiyang
* @Time 2022/10/14 15:17
*/
class MenuContentItem extends BaseModel {
  MenuContentItem.from(obj);

  @override
  Map<String, dynamic> to() => {};
}
