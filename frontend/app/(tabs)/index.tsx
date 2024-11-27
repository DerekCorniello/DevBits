import { Post } from "@/components/Post";
import ScrollView from "@/components/ScrollView";
import CreatePost from "@/components/CreatePost";
import TopBar from "@/components/ui/TopBar";

export default function HomeScreen() {
  return (
    <>
      {/* <TopBar /> */}
      <ScrollView>
        <Post
          ID={1}
          User={"Username"}
          Project={1}
          Likes={69}
          Content="This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app."
          CreationDate={new Date()}
          Comments={[]}
        />
      </ScrollView>
      <CreatePost />
    </>
  );
}
