import * as THREE from 'three';
import { TransformControls } from 'three/addons/controls/TransformControls.js';

export class Vector3ParameterNodeController {
    constructor(lightNode, nodeManager, id, parameterData, app) {
        this.nodeManager = nodeManager;
        this.id = id;
        this.updating = false;

        const control = new TransformControls(app.Camera, app.Renderer.domElement);
        control.setMode('translate');
        control.setSpace("local");

        this.mesh = new THREE.Group();

        control.addEventListener('dragging-changed', (event) => {
            app.OrbitControls.enabled = !event.value;

            if (!app.OrbitControls.enabled) {
                return;
            }

            nodeManager.nodeParameterChanged({
                id: id,
                data: {
                    x: this.mesh.position.x,
                    y: this.mesh.position.y,
                    z: this.mesh.position.z,
                },
                binary: false
            });
        });

        app.ViewerScene.add(this.mesh);

        this.lightNode = lightNode;

        const curVal = parameterData.currentValue;
        this.lightNode.setProperty("x", curVal.x);
        this.lightNode.setProperty("y", curVal.y);
        this.lightNode.setProperty("z", curVal.z);
        this.mesh.position.set(curVal.x, curVal.y, curVal.z);

        const helper = control.getHelper();
        app.Scene.add(helper)
        control.attach(this.mesh);

        this.lightNode.setTitle(parameterData.name);

        helper.visible = false;
        helper.enabled = false;
        control.enabled = false;

        this.lightNode.addSelectListener(() => {
            helper.visible = true;
            helper.enabled = true;
            control.enabled = true;
        });

        this.lightNode.addUnselectListener(() => {
            helper.visible = false;
            helper.enabled = false;
            control.enabled = false;
        });

        this.lightNode.addPropertyChangeListener("x", this.propertyChange.bind(this));
        this.lightNode.addPropertyChangeListener("y", this.propertyChange.bind(this));
        this.lightNode.addPropertyChangeListener("z", this.propertyChange.bind(this));
    }

    propertyChange() {
        if (this.updating) {
            return
        }
        this.nodeManager.nodeParameterChanged({
            id: this.id,
            data: {
                x: this.lightNode.getProperty("x"),
                y: this.lightNode.getProperty("y"),
                z: this.lightNode.getProperty("z")
            },
            binary: false
        });
    }

    update(parameterData) {
        this.updating = true;
        const curVal = parameterData.currentValue;
        this.mesh.position.set(curVal.x, curVal.y, curVal.z);
        this.lightNode.setProperty("x", curVal.x);
        this.lightNode.setProperty("y", curVal.y);
        this.lightNode.setProperty("z", curVal.z);
        this.updating = false;
    }
}     