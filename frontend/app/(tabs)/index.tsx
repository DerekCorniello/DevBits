import { Post } from "@/components/Post";
import ScrollView from "@/components/ScrollView";
import CreatePost from "@/components/CreatePost";
import TopBar from "@/components/ui/TopBar";
import { PostProps } from "@/constants/Types";

export default function HomeScreen() {
  return (
    <>
      {/* <TopBar /> */}
      <ScrollView>
        <Post
          ID={1}
          User={2}
          Project={1}
          Likes={69}
          Content="This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app."
          CreationDate="2021-01-01T00:00:00Z"
          Comments={[]}
        />
      </ScrollView>
      <CreatePost />
    </>
  );
}
