import QtQuick 2.0
import QtQuick.Controls 1.4

Row {
    property string desc
    property alias value: slider.value

    width: parent.width
    Label {
        text: desc
        width: parent.width / 3
    }

    Slider {
        id: slider
        minimumValue: 0
        maximumValue: 1
        stepSize: 0.001
        width: parent.width / 3

    }

    Label {
        // spacer
        text: "  "
    }

    Label {
        text: parent.value
        width: parent.width / 3
    }
}
