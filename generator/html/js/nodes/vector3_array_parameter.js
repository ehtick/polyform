import { TransformControls } from 'three/addons/controls/TransformControls.js';
import * as THREE from 'three';


export class Vector3ArrayParameterNodeController {
    constructor(lightNode, nodeManager, id, parameterData, app) {
        this.id = id;
        this.nodeManager = nodeManager;
        this.app = app;
        this.scene = app.ViewerScene;
        this.allPositionControlHelpers = [];
        this.allPositionControls = [];
        this.allPositionControlHelpersMeshes = [];
        this.renderControls = false;

        parameterData.currentValue?.forEach((ele) => {
            this.newPositionControl(ele);
        })

        this.lightNode = lightNode;
        this.lightNode.setTitle(parameterData.name);
        const addPointButton = GlobalWidgetFactory.create(this.lightNode, "button", {
            text: "Add Point",
            callback: () => {
                const paramData = this.buildParameterData();

                if (paramData.length > 0) {
                    const oldEle = paramData[paramData.length - 1]
                    const newEle = {
                        x: oldEle.x + 1,
                        y: oldEle.y,
                        z: oldEle.z,
                    }

                    paramData.push(newEle)
                } else {
                    paramData.push({ x: 0, y: 0, z: 0 })
                }


                this.nodeManager.nodeParameterChanged({
                    id: this.id,
                    data: paramData,
                });
            }
        })

        this.lightNode.addWidget(addPointButton);


        this.lightNode.addSelectListener((obj) => {
            this.renderControls = true;
            this.updateControlRendering();
        });

        this.lightNode.addUnselectListener((obj) => {
            this.renderControls = false;
            this.updateControlRendering();
        });
    }

    updateControlRendering() {
        this.allPositionControlHelpers.forEach((v) => {
            v.visible = this.renderControls;
            v.enabled = this.renderControls;
        });
        this.allPositionControls.forEach((v) => {
            v.enabled = this.renderControls;
        });
    }

    buildParameterData() {
        const data = [];

        this.allPositionControlHelpersMeshes.forEach((ele) => {
            data.push({
                x: ele.position.x,
                y: ele.position.y,
                z: ele.position.z,
            })
        })

        return data
    }

    newPositionControl(pos) {
        const control = new TransformControls(this.app.Camera, this.app.Renderer.domElement);
        control.setMode('translate');
        control.setSpace("local");


        const mesh = new THREE.Group();

        control.addEventListener('dragging-changed', (event) => {
            this.app.OrbitControls.enabled = !event.value;

            if (this.app.OrbitControls.enabled) {
                this.nodeManager.nodeParameterChanged({
                    id: this.id,
                    data: this.buildParameterData()
                });
            }
        });

        const helper = control.getHelper();
        this.allPositionControlHelpers.push(helper);
        this.allPositionControls.push(control);
        this.allPositionControlHelpersMeshes.push(mesh);

        this.scene.add(mesh);
        this.app.Scene.add(helper);
        mesh.position.set(pos.x, pos.y, pos.z);
        control.attach(mesh);

        helper.visible = this.renderControls;
        helper.enabled = this.renderControls;
        control.enabled = this.renderControls;
    }

    clearPositionControls() {
        this.allPositionControlHelpers.forEach((v) => {
            this.app.Scene.remove(v);
        });
        this.allPositionControlHelpersMeshes.forEach((v) => {
            this.scene.remove(v);
        })
        this.allPositionControlHelpers = [];
        this.allPositionControlHelpersMeshes = [];
        this.allPositionControls = [];
    }

    update(parameterData) {
        this.clearPositionControls();
        parameterData.currentValue?.forEach((ele) => {
            this.newPositionControl(ele);
        })
    }
}