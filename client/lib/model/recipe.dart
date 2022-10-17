import 'package:client/common/model.dart';
import 'package:client/model/activity.dart';
import 'package:client/model/model.dart';
import 'package:client/model/post.dart';
import 'package:client/model/tag.dart';

/*
* 食谱信息
* @author wuxubaiyang
* @Time 2022/10/14 15:16
*/
class RecipeModel extends BaseModel with BasePart, CreatorPart {
  // 食谱标题
  final String title;

  // 食谱描述
  final String desc;

  // 食谱图片集合
  final List<String> images;

  // 预计耗时
  final num time;

  // 难度评分
  final num rating;

  // 步骤集合
  final List<RecipeStepItem> steps;

  // 菜系字典码集合
  final List<String> cuisineCodes;

  // 口味字典码集合
  final List<String> tasteCodes;

  // 主料字典码集合
  final List<String> ingredientsMainCodes;

  // 辅料字典码集合
  final List<String> ingredientsSubCodes;

  // 标签字典码集合
  final List<String> tagCodes;

  // 标签集合
  final List<TagModel> tags;

  // 活动记录id
  final String? activityRecordId;

  // 活动记录信息
  late ActivityRecordModel? activityRecord;

  // 是否点赞
  final bool liked;

  // 点赞数量
  final num likeCount;

  // 是否收藏
  final bool collected;

  // 收藏数量
  final num collectCount;

  RecipeModel.from(obj)
      : title = obj?["title"] ?? "",
        desc = obj?["desc"] ?? "",
        images = obj?["images"] ?? [],
        time = obj?["time"] ?? 0,
        rating = obj?["rating"] ?? 0,
        steps = (obj?["steps"] ?? [])
            .map<RecipeStepItem>((e) => RecipeStepItem.from(e))
            .toList(),
        cuisineCodes =
            (obj?["cuisineCodes"] ?? []).map<String>((e) => "$e").toList(),
        tasteCodes =
            (obj?["tasteCodes"] ?? []).map<String>((e) => "$e").toList(),
        ingredientsMainCodes = (obj?["ingredientsMainCodes"] ?? [])
            .map<String>((e) => "$e")
            .toList(),
        ingredientsSubCodes = (obj?["ingredientsSubCodes"] ?? [])
            .map<String>((e) => "$e")
            .toList(),
        tagCodes = (obj?["tagCodes"] ?? []).map<String>((e) => "$e").toList(),
        tags = (obj?["tags"] ?? [])
            .map<TagModel>((e) => TagModel.from(e))
            .toList(),
        activityRecordId = obj?["activityRecordId"],
        liked = obj?["liked"] ?? false,
        likeCount = obj?["likeCount"] ?? 0,
        collected = obj?["collected"] ?? false,
        collectCount = obj?["collectCount"] ?? 0 {
    initBasePart(obj);
    initCreatorPart(obj);
    if (obj?["activityRecord"] != null) {
      activityRecord = ActivityRecordModel.from(obj?["activityRecord"] ?? {});
    }
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        ...creatorPart,
        "title": title,
        "desc": desc,
        "images": images,
        "time": time,
        "rating": rating,
        "steps": steps.map((e) => e.to()).toList(),
        "cuisineCodes": cuisineCodes,
        "tasteCodes": tasteCodes,
        "ingredientsMainCodes": ingredientsMainCodes,
        "ingredientsSubCodes": ingredientsSubCodes,
        "tagCodes": tagCodes,
        "tags": tags.map((e) => e.to()).toList(),
        "activityRecordId": activityRecordId,
        "activityRecord": activityRecord?.to(),
        "liked": liked,
        "likeCount": likeCount,
        "collected": collected,
        "collectCount": collectCount,
      };

  // 编辑结构体
  Map<String, dynamic> toUpdateInfo() => {
        "title": title,
        "desc": desc,
        "images": images,
        "time": time,
        "rating": rating,
        "steps": steps.map((e) => e.to()).toList(),
        "cuisineCodes": cuisineCodes,
        "tasteCodes": tasteCodes,
        "ingredientsMainCodes": ingredientsMainCodes,
        "ingredientsSubCodes": ingredientsSubCodes,
        "tagCodes": tagCodes,
        "activityRecordId": activityRecordId,
      };
}

/*
* 食谱步骤项
* @author wuxubaiyang
* @Time 2022/10/14 15:17
*/
class RecipeStepItem extends BaseModel {
  // 预计耗时
  final num time;

  // 步骤内容
  final List<ContentItem> contents;

  RecipeStepItem.from(obj)
      : time = obj?["time"] ?? 0,
        contents = (obj?["contents"] ?? [])
            .map<ContentItem>((e) => ContentItem.from(e))
            .toList();

  @override
  Map<String, dynamic> to() => {
        "time": time,
        "contents": contents.map((e) => e.to()).toList(),
      };
}
