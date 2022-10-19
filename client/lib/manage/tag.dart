import 'dart:convert';

import 'package:client/common/manage.dart';
import 'package:client/model/tag.dart';
import 'package:flutter/cupertino.dart';
import 'package:json_path/json_path.dart';

/*
* 标签管理
* @author wuxubaiyang
* @Time 2022/3/17 14:14
*/
class TagManage extends BaseManage {
  static final TagManage _instance = TagManage._internal();

  factory TagManage() => _instance;

  TagManage._internal();

  // 定义字典检索对象
  final jsonPath = JsonPath(r'$..[?hasCode]');

  // 根据code集合查找目标标签对象
  Future<List<TagModel>> findTags(
    BuildContext context, {
    required TagSource source,
    required List<String> codes,
  }) async {
    var json = await loadTagsMapSource(context, source);
    return jsonPath
        .read(json, filters: {
          "hasCode": (m) => m.value is Map && codes.contains(m.value["code"]),
        })
        .map((e) => TagModel.from(e.value))
        .toList();
  }

  // 根据code查找一个标签
  Future<TagModel?> findTag(
    BuildContext context, {
    required TagSource source,
    required String code,
  }) async {
    var result = await findTags(
      context,
      source: source,
      codes: [code],
    );
    return result.isNotEmpty ? result.first : null;
  }

  // 加载某个标签对象
  Future<List> loadTagsMapSource(
    BuildContext context,
    TagSource source,
  ) async {
    var raw = await _loadAssetFileAsString(context, source.path);
    return jsonDecode(raw);
  }

  // 缓存assetsBundle对象
  AssetBundle? assetBundle;

  // 加载assets资源
  Future<String> _loadAssetFileAsString(BuildContext context, String path) {
    assetBundle ??= DefaultAssetBundle.of(context);
    return assetBundle!.loadString(path);
  }
}

// 单例调用
final tagManage = TagManage();

// 标签来源枚举
enum TagSource {
  activityType,
  address,
  commentType,
  medalRarity,
  noticeType,
  recipeCuisine,
  recipeIngredientsMain,
  recipeIngredientsSub,
  recipeTaste,
  userEvaluate,
  userGender,
}

// 标签根路径
const String tagRoot = "assets/tags";

// 标签来源枚举扩展
extension TagSourceExtension on TagSource {
  // 获取资源路径
  String get path => <TagSource, String>{
        TagSource.activityType: "$tagRoot/activity_type.json",
        TagSource.address: "$tagRoot/address.json",
        TagSource.commentType: "$tagRoot/comment_type.json",
        TagSource.medalRarity: "$tagRoot/medal_rarity.json",
        TagSource.noticeType: "$tagRoot/notice_type.json",
        TagSource.recipeCuisine: "$tagRoot/recipe_cuisine.json",
        TagSource.recipeIngredientsMain:
            "$tagRoot/recipe_ingredients_main.json",
        TagSource.recipeIngredientsSub: "$tagRoot/recipe_ingredients_sub.json",
        TagSource.recipeTaste: "$tagRoot/recipe_taste.json",
        TagSource.userEvaluate: "$tagRoot/user_evaluate.json",
        TagSource.userGender: "$tagRoot/user_gender.json",
      }[this]!;
}
