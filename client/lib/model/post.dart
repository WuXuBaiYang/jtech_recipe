import 'package:client/common/model.dart';
import 'package:client/model/activity.dart';
import 'package:client/model/model.dart';
import 'package:client/model/recipe.dart';
import 'package:client/model/tag.dart';

/*
* 帖子信息
* @author wuxubaiyang
* @Time 2022/10/14 15:16
*/
class PostModel extends BaseModel with BasePart, CreatorPart {
  // 帖子标题
  final String title;

  // 内容集合
  final List<ContentItem> contents;

  // 标签代码集合
  final List<String> tagCodes;

  // 标签集合
  final List<TagModel> tags;

  // 活动记录id
  final String? activityRecordId;

  // 活动记录信息
  late ActivityRecordModel? activityRecord;

  // 食谱id
  final String? recipeId;

  // 食谱信息
  late RecipeModel? recipe;

  // 是否点赞
  final bool liked;

  // 点赞数量
  final num likeCount;

  // 是否收藏
  final bool collected;

  // 收藏数量
  final num collectCount;

  PostModel.from(obj)
      : title = obj?["title"] ?? "",
        contents = (obj?["contents"] ?? [])
            .map<ContentItem>((e) => ContentItem.from(e))
            .toList(),
        tagCodes = (obj?["tagCodes"] ?? []).map<String>((e) => "$e").toList(),
        tags = (obj?["tags"] ?? [])
            .map<TagModel>((e) => TagModel.from(e))
            .toList(),
        activityRecordId = obj?["activityRecordId"],
        recipeId = obj?["recipeId"],
        liked = obj?["liked"] ?? false,
        likeCount = obj?["likeCount"] ?? 0,
        collected = obj?["collected"] ?? false,
        collectCount = obj?["collectCount"] ?? 0 {
    initBasePart(obj);
    initCreatorPart(obj);
    if (obj?["activityRecord"] != null) {
      activityRecord = ActivityRecordModel.from(obj?["activityRecord"] ?? {});
    }
    if (obj?["recipe"] != null) {
      recipe = RecipeModel.from(obj?["recipe"] ?? {});
    }
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        ...creatorPart,
        "title": title,
        "contents": contents.map((e) => e.to()).toList(),
        "tagCodes": tagCodes,
        "tags": tags.map((e) => e.to()).toList(),
        "activityRecordId": activityRecordId,
        "activityRecord": activityRecord?.to(),
        "recipeId": recipeId,
        "recipe": recipe?.to(),
        "liked": liked,
        "likeCount": likeCount,
        "collected": collected,
        "collectCount": collectCount,
      };

  // 整理更新结构
  Map<String, dynamic> toUpdateInfo() => {
        "title": title,
        "contents": contents.map((e) => e.to()).toList(),
        "tagCodes": tagCodes,
        "activityRecordId": activityRecordId,
        "recipeId": recipeId,
      };
}

/*
* 帖子内容结构
* @author wuxubaiyang
* @Time 2022/10/14 15:17
*/
class ContentItem extends BaseModel {
  ContentItem.from(obj);

  @override
  Map<String, dynamic> to() => {};
}
