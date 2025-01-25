
import { useState, useRef } from "react";
import { Overlay } from '@rneui/themed';
import { View, Text, StyleSheet, Animated } from "react-native";
import { Card, Icon, CheckBox } from "@rneui/themed";
import { ThemedText } from "@/components/ThemedText";
import { useThemeColor } from "@/hooks/useThemeColor";
import { CommentProps } from "@/constants/Types";

export function Comment({
    id,
    user,
    likes,
    parent_comment,
    created_on,
    content,
}: PostProps) {
    const cardBackgroundColor = useThemeColor(
        { light: "light grey", dark: "#151515" },
        "background"
    );
    let CreationDate = new Date(created_on);
    return (
        
    );
}
