import { FAB } from "@rneui/themed";
import { Icon } from "@rneui/themed";

export default function CreatePost() {
  return (
    <FAB
      visible={true}
      placement="right"
      icon={<Icon name="code" type="octicon" color="black" size={25} />}
      color="#16ff00"
      style={{ marginBottom: 100 }}
      size="large"
    />
  );
}
