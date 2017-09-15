import QtQuick 2.0
import QtQuick.Controls 1.4

ApplicationWindow {

    width: 800
    height: 600
    title: "Simple"

    Column {
        width: parent.width
        Parameter {
            id: p
            desc: "Fraction of honest nodes (p)"
        }
        Parameter {
            id: n
            desc: "Number of subtasks per node (n)"
            step: 1
            minValue: 1
            maxValue: 100
        }
        Parameter {
            id: c
            desc: "Fraction of correct subtasks (c)"
        }

        Label {
            text: "Expected voting result"
        }
        Row {
            Label {
                text: "Scenario: catching dishonest peer: "
            }
            Label {
                property double expected: (1 - p.value) + p.value * Math.pow(c.value, n.value)
                text: expected
                color: (expected < 0.5) ? "black" : "red"
            }
        }
        Row {
            Label {
                text: "Scenario: validating a honest peer: "
            }
            Label {
                text: p.value
                color: (p.value > 0.5) ? "black" : "red"
            }
        }
    }
}
