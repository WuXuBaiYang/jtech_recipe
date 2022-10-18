import 'package:client/common/manage.dart';

/*
* 标签管理
* @author wuxubaiyang
* @Time 2022/3/17 14:14
*/
class TagManage extends BaseManage {
  static final TagManage _instance = TagManage._internal();

  factory TagManage() => _instance;

  TagManage._internal();

// var j = JsonPath(r'$..[?hasCode]');
// var raw = await DefaultAssetBundle.of(context)
//     .loadString("assets/tags/address.json");
// var result = j
//     .read(jsonDecode(raw), filters: {
//   "hasCode": (m) {
//     return m.value is Map &&
//         ["110102", "120105"].contains(m.value["code"]);
//   },
// })
//     .map((e) => TagModel.from(e.value))
//     .toList();
// print("object");
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
