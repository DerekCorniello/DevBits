import { Post } from "@/components/Post";
import ScrollView from "@/components/ScrollView";

export default function HomeScreen() {
  return (
    <ScrollView>
      <Post
        ID={1}
        User={1}
        Project={1}
        Likes={1}
        Content="This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app."
        Comments={[]}
        CreationDate={new Date()}
      />
    </ScrollView>
  );
}
