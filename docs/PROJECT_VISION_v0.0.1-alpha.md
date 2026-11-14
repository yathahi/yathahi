# Yatahi — Project Vision

**Document version:** 0.0.1-alpha  
**Document type:** Project vision

## 1. What is Yatahi?

Yatahi is an ecosystem for creating, connecting, securing, automating and maintaining infrastructure as one coherent environment.

That environment can be as simple as a single computer running several services, a small home laboratory with two devices, or a collection of systems spread across many locations. Yatahi is intended to make sense at every point along that range. It should not require a large network, a dedicated infrastructure team or an enterprise-sized deployment to be useful.

Yatahi brings the devices, systems and services that make up an environment together so they can be understood and managed in relation to one another. As the environment grows, more devices and services can become part of it without changing the overall purpose: to make the infrastructure easier to create, operate, secure and maintain.

People should be able to focus on what they want their systems to do, rather than needing to become experts in every technology involved. Yatahi aims to take on much of the complexity behind turning those needs into working infrastructure.

Yatahi is not limited to web hosting. Hosting is an initial area of focus, while the wider vision includes many kinds of systems, services and devices.

## 2. Why Yatahi exists

Even a small setup can involve several services, configurations and technical tasks. As more services or devices are added, keeping everything working together can become difficult. The same challenge exists in a home laboratory, a small office, a research environment or a much larger deployment; the scale changes, but the need to create, coordinate, protect and maintain the systems remains.

Yatahi aims to make that work more approachable by bringing these responsibilities into one ecosystem. It is not just about making one service easier to install. It is about helping people manage the environment around their services as a whole, without having to handle every underlying detail themselves.

Whether someone is maintaining a personal server or looking after infrastructure shared by a large team, they should be able to spend more time using and improving their systems and less time dealing with repetitive or unnecessarily complicated infrastructure work.

## 3. From one computer to infrastructure of any scale

Yatahi should be useful even when there is only one device. A single computer running Yatahi and hosting multiple services is a valid and meaningful setup. There does not need to be another machine, a separate server, or a distributed network for Yatahi to have a purpose.

A setup with two devices should be just as welcome. So should a home lab with several machines, a small office, a research laboratory, or an infrastructure spanning many locations and containing millions of devices. None of these environments should be considered outside Yatahi's intended scope because they are too small or too large.

The goal is not to make every setup behave like a large deployment. A simple environment should be able to remain simple, without being burdened by complexity it does not need. At the same time, a growing environment should be able to add devices, services and locations without outgrowing the ecosystem. Yatahi should adapt to the scale and needs of the environment, rather than forcing every user into the same model.

## 4. One environment across devices and locations

When an environment includes multiple devices, Yatahi aims to bring their participating devices and services together as one logical infrastructure. Those devices may be beside one another, in different buildings, or in different parts of the world. Physical distance should not prevent them from belonging to the same environment.

Not every device needs to host a service. Some devices may simply participate as members of the environment, while others provide hosting, DNS, backend systems or specialised services. Different devices can have different purposes while remaining part of the same wider setup.

For example, a personal home server might begin by hosting several services on one computer. Its owner could later add another machine, move a service to it, or include devices at another location. In a larger setting, the same broad idea can connect systems used by different teams or sites. Yatahi's vision is to support that progression without treating a single computer, a small local network and a widely distributed infrastructure as fundamentally unrelated use cases.

## 5. What Yatahi aims to do

Yatahi has four overarching objectives.

### 5.1 Infrastructure creation

Yatahi aims to help create the infrastructure needed for a particular purpose, using the devices and resources available. This may mean setting up several services on one computer, distributing services across a home lab, or establishing a much larger environment.

Creating infrastructure should not require users to understand or manually perform every underlying technical task. The intended experience is to describe the outcome they need and have Yatahi handle as much of the work involved in establishing it as possible.

### 5.2 Infrastructure orchestration

Yatahi aims to coordinate the devices, services and resources that make up an environment, so they can work together as a coherent whole. This applies whether everything runs on one machine or is distributed across many devices and locations.

People should be able to understand what is present, how its parts relate and how changes affect the wider environment, rather than having to manage each service or device in isolation. As a setup grows or changes, Yatahi should help keep it understandable and coordinated.

For environments shared by multiple people, the vision includes the ability to have multiple Administrators with equal authority. Responsibility should not have to depend on one particular person or one physical device.

### 5.3 Infrastructure security

Security is a fundamental part of Yatahi's purpose, not an optional feature added after infrastructure has been created.

Yatahi aims to provide coordinated protection across the environment, including its network, individual devices, and the services and workloads running on them. This protection should complement the security capabilities already provided by the underlying systems and services, not replace them.

A problem or compromise affecting one part of an environment should not automatically put everything else at risk. Yatahi aims to help identify threats, contain and isolate affected parts, and protect the rest of the environment. These principles matter in a small setup just as they do in a large one; the level of complexity may differ, but security should remain part of the design.

### 5.4 Infrastructure maintenance

Yatahi aims to help keep infrastructure healthy and dependable after it has been created. This includes monitoring the environment, recognising unexpected conditions, diagnosing problems, taking corrective action and supporting recovery.

A person maintaining a home server should not have to manually notice every service failure or developing issue. The same applies to people responsible for larger environments. Yatahi's long-term aim is to recognise and address problems automatically wherever it can do so safely, while making important conditions and actions visible to the people responsible.

## 6. Accessible by design

Yatahi aims to make infrastructure management approachable to people with different levels of technical experience. Someone running a personal server should not need the same specialist knowledge expected of a professional systems administrator just to manage basic services. Likewise, larger and more complex environments should not require every routine task to be handled manually by specialists.

Complexity should be handled by the ecosystem wherever possible, rather than passed on to the user. People should be able to express what they want to achieve without having to know every detail of the technologies involved. Human supervision and control should remain available, especially where judgement, approval or safety matters, but deep technical knowledge should not be a prerequisite for using Yatahi.

The intended experience includes a clear visual understanding of the environment: its devices and services, their condition, and how they are connected. People should be able to understand and adjust their setup through an approachable visual interface. Where a proposed change could cause a problem, Yatahi should explain the concern and help the user make an informed choice.

## 7. Infrastructure that remains independent

Yatahi is intended to create and improve infrastructure, not make that infrastructure permanently dependent on Yatahi's presence.

Once established, the underlying devices, services and configurations should continue to operate independently. If Yatahi becomes unavailable or is removed, the infrastructure it helped create should remain functional. The user would lose Yatahi's coordinated management, automation, additional security capabilities, monitoring, diagnosis and automated maintenance, but the underlying services should not stop working merely because Yatahi is no longer present.

The underlying systems retain their own capabilities and protections. Yatahi's role is to add coordination, management and security around them, not to become a mandatory runtime dependency for every service it helps establish.

## 8. Guiding principles

Yatahi's vision is guided by the following principles:

- **Useful at any scale:** Work meaningfully on one computer running multiple services, in a two-device home lab, in a small office, and across infrastructures containing millions of devices.
- **No unnecessary complexity:** Keep simple setups simple, while allowing more involved environments to grow without outgrowing the ecosystem.
- **Unified environment:** Bring participating devices and services together as parts of one coherent infrastructure, whether they are local or far apart.
- **Automation over manual complexity:** Reduce the need for users to perform specialist, repetitive or complicated infrastructure work.
- **Security throughout:** Make protection a fundamental part of infrastructure creation and operation, across networks, devices and workloads.
- **Isolation and resilience:** Limit the impact that a failure or compromise in one part can have on the rest of the environment.
- **Operational continuity:** Preserve the independent operation of established infrastructure, including when Yatahi is unavailable or removed.
- **Broad compatibility:** Aim to support the varied devices, services and needs found in personal, small-scale and large-scale environments.
- **Human understanding and oversight:** Give people visibility and meaningful control without requiring them to master every technical detail.
- **Continuous improvement:** Help infrastructure remain healthy, secure and adaptable as needs change.

## 9. A vision that can evolve

This document records the vision of Yatahi at version **0.0.1-alpha**. It describes the project's current purpose, ambitions and guiding principles; it is not a permanent limitation on what Yatahi may become.

As the project develops and its understanding of infrastructure evolves, this vision may also change. Future versions may refine, expand or redefine its objectives. The version history of this document should make that evolution visible while preserving the context of earlier visions.

The enduring intention is to make creating and operating infrastructure of any scale simpler, more secure, more coordinated and more resilient—whether that infrastructure lives on one computer or spans many locations—while allowing the infrastructure itself to remain independent.
