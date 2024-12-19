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
