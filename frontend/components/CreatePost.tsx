import { FAB } from "@rneui/themed";
import { Icon } from "@rneui/themed";

export default function CreatePost() {
  return (
    <FAB
      visible={true}
      placement="right"
      icon={<Icon name="code" type="octicon" color="white" size={25} />}
      color="#004477"
      style={{ marginBottom: 100 }}
      size="large"
    />
  );
}
