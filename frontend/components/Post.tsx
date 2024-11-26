import { CheckBox, Icon } from "@rneui/themed";
import { Card } from "@rneui/themed";
import { ThemedText } from "@/components/ThemedText";
import { useState } from "react";
import { useThemeColor } from "@/hooks/useThemeColor";
import { Stack } from "expo-router";
import { ThemedView } from "./ThemedView";

function Like() {
  const [checked, setState] = useState(true);
  const toggleCheckbox = () => setState(!checked);
  return (
    <CheckBox
      center
      checkedIcon={
        <Icon
          name="heart-fill"
          type="octicon"
          color="green"
          size={25}
          iconStyle={{ marginRight: 10 }}
        />
      }
      uncheckedIcon={
        <Icon
          name="heart"
          type="octicon"
          color="grey"
          size={25}
          iconStyle={{ marginRight: 10 }}
        />
      }
      checked={checked}
      onPress={() => toggleCheckbox()}
    />
  );
}

export type PostProps = {
  ID: number;
  User: number; // Linked to User
  Project: number; // Linked to Project
  Likes: number;
  Content: string;
  Comments: number[];
  CreationDate: Date;
};

export function Post({
  ID,
  User,
  Project,
  Likes,
  Content,
  Comments,
  CreationDate,
}: PostProps) {
  return (
    <ThemedView>
      <Card>
        <Card.Title>{User}</Card.Title>
        <Card.Divider />
        <ThemedText type="defaultSemiBold">{Content}</ThemedText>
        <Card.Divider />
        <Like />
        <ThemedText type="subtitle">{Likes}</ThemedText>
        <ThemedText type="subtitle">{CreationDate.toUTCString()}</ThemedText>
      </Card>
    </ThemedView>
  );
}
