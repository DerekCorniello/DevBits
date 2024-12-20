import { Post } from "@/components/Post";
import { useThemeColor } from "@/hooks/useThemeColor";
import { StyleSheet, ScrollView, View } from "react-native";
import CreatePost from "@/components/CreatePost";
import { MyFilter } from "@/components/filter";

export default function HomeScreen() {
  const cardBackgroundColor = useThemeColor(
    { light: "#fff", dark: "#040607" },
    "background"
  );
  return (
    <>
      <View style={styles.filterContainer}>
        <MyFilter />
      </View>
      <ScrollView>
        <Post
          id={1}
          user={2}
          project={1}
          likes={69}
          content="This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app."
          created_on="2021-01-01T00:00:00Z"
          comments={[]}
        />
      </ScrollView>
      <CreatePost />
    </>
  );
}

const styles = StyleSheet.create({
  filterContainer: {
    position: "relative",
    alignSelf: "flex-start",
    padding: 20,
    zIndex: 1,
  },
});
